package queue

import (
	"context"
	"sync"
	"time"

	"github.com/celsiainternet/elvis/logs"
)

/**
* Queue: Cola que despacha lotes al llegar a maxEvents o al vencer period.
**/
type Queue[T any] struct {
	queue     chan T
	handler   func(context.Context, []T) error
	maxEvents int
	period    time.Duration
	mu        sync.RWMutex
	closed    bool
	done      chan struct{}
}

/**
* New: Crea una cola que despacha lotes de maxEvents o cada period.
* @param queueSize int, tamaño máximo de la cola; maxEvents int, cantidad máxima de eventos por lote; period time.Duration, tiempo máximo de espera por lote; handler func(context.Context, []T) error
* @return *Queue[T]
**/
func New[T any](
	queueSize int,
	maxEvents int,
	period time.Duration,
	handler func(context.Context, []T) error,
) *Queue[T] {
	if maxEvents <= 0 {
		maxEvents = 1
	}

	if period <= 0 {
		period = time.Second
	}

	q := &Queue[T]{
		queue:     make(chan T, queueSize),
		handler:   handler,
		maxEvents: maxEvents,
		period:    period,
		done:      make(chan struct{}),
	}

	go q.worker()

	return q
}

/**
* worker: Acumula items y despacha el lote por tamaño o por tiempo.
**/
func (q *Queue[T]) worker() {
	defer close(q.done)

	ctx := context.Background()

	buffer := make([]T, 0, q.maxEvents)

	timer := time.NewTimer(q.period)
	defer timer.Stop()

	dispatch := func(batch []T) {
		if err := q.handler(ctx, batch); err != nil {
			logs.Alert(err)
		}
	}

	resetTimer := func() {
		if !timer.Stop() {
			select {
			case <-timer.C:
			default:
			}
		}
		timer.Reset(q.period)
	}

	for {
		select {
		case item, ok := <-q.queue:
			if !ok {
				if len(buffer) > 0 {
					dispatch(buffer)
				}
				return
			}

			buffer = append(buffer, item)

			// Condición de disparo 1: tamaño de lote alcanzado.
			if len(buffer) >= q.maxEvents {
				batch := buffer
				buffer = make([]T, 0, q.maxEvents)
				resetTimer()
				dispatch(batch)
			}
		case <-timer.C:
			// Condición de disparo 2: venció el tiempo de espera.
			if len(buffer) > 0 {
				batch := buffer
				buffer = make([]T, 0, q.maxEvents)
				dispatch(batch)
			}
			timer.Reset(q.period)
		}
	}
}

/**
* Push: Encola un item; no hace nada si la cola está cerrada.
* @param item T
**/
func (q *Queue[T]) Push(item T) {
	q.mu.RLock()
	defer q.mu.RUnlock()

	if q.closed {
		return
	}

	q.queue <- item
}

/**
* Len: Retorna la cantidad de items pendientes en la cola.
* @return int
**/
func (q *Queue[T]) Len() int {
	return len(q.queue)
}

/**
* Close: Cierra la cola y espera el último lote.
**/
func (q *Queue[T]) Close() {
	q.mu.Lock()
	if q.closed {
		q.mu.Unlock()
		return
	}

	q.closed = true
	close(q.queue)
	q.mu.Unlock()

	<-q.done
}
