package service

import (
	"errors"
	"fmt"
	"net"
	"net/rpc"
	"testing"

	"github.com/celsiainternet/elvis/et"
	"github.com/celsiainternet/elvis/jrpc"
)

type fakeCall struct {
	method string
	args   et.Json
}

/**
* withFakes: Reemplaza solver, llamada y event.Work; devuelve los registros de cada uno.
* @param t *testing.T, solveErr error, item et.Item, callErr error
* @return *[]string, *[]fakeCall, *[]fakeCall
**/
func withFakes(t *testing.T, solveErr error, item et.Item, callErr error) (*[]string, *[]fakeCall, *[]fakeCall) {
	t.Helper()
	solved := []string{}
	calls := []fakeCall{}
	works := []fakeCall{}
	oldSolve, oldCall, oldWork := rpcSolve, rpcCall, workPublish
	t.Cleanup(func() { rpcSolve, rpcCall, workPublish = oldSolve, oldCall, oldWork })

	rpcSolve = func(method string) (*jrpc.Solver, error) {
		solved = append(solved, method)
		if solveErr != nil {
			return nil, solveErr
		}
		return &jrpc.Solver{Method: method, Host: "fake", Port: 1}, nil
	}
	rpcCall = func(solver *jrpc.Solver, args et.Json) (et.Item, error) {
		calls = append(calls, fakeCall{method: solver.Method, args: args})
		return item, callErr
	}
	workPublish = func(event string, data et.Json) et.Json {
		works = append(works, fakeCall{method: event, args: data})
		return et.Json{"_id": data.Str("service_id"), "event": event, "data": data}
	}

	return &solved, &calls, &works
}

func TestSendRpcMethods(t *testing.T) {
	t.Setenv("SEND_RPC", "")
	t.Setenv("OTP_RPC", "")
	cases := []struct {
		name   string
		method string
		run    func() et.Json
	}{
		{"sms", "Service-send.Services.SendSMS", func() et.Json {
			return SendSms("p1", "s1", []string{"3000000000"}, "hola", nil, TpTransactional, et.Json{})
		}},
		{"whatsapp", "Service-send.Services.SendWhatsapp", func() et.Json {
			return SendWhatsapp("p1", "s1", 21, []string{"3000000000"}, nil, TpTransactional, et.Json{})
		}},
		{"email", "Service-send.Services.SendEmail", func() et.Json {
			return SendEmail("p1", "s1", []et.Json{{"email": "a@b.co"}}, "asunto", "<p/>", nil, TpTransactional, et.Json{})
		}},
		{"emailByTemplate", "Service-send.Services.SendEmailByTemplate", func() et.Json {
			return SendEmailByTemplate("p1", "s1", []string{"g1"}, "asunto", "t1", nil, TpTransactional, et.Json{})
		}},
		{"emailToByTemplate", "Service-send.Services.SendEmailToByTemplate", func() et.Json {
			return SendEmailToByTemplate("p1", "s1", "Ana", "a@b.co", "asunto", "t1", nil, TpTransactional, et.Json{})
		}},
		{"otpSms", "Service-otp.Services.GenerateOTPSMS", func() et.Json {
			return SendOtpSms("p1", "3000000000", "Ana", "app", 6, 5, et.Json{})
		}},
		{"otpWhatsapp", "Service-otp.Services.GenerateOTPWhatsapp", func() et.Json {
			return SendOtpWhatsapp("p1", "3000000000", "Ana", "app", 6, 5, et.Json{})
		}},
		{"otpEmail", "Service-otp.Services.GenerateOTPEmail", func() et.Json {
			return SendOtpEmail("p1", "a@b.co", "Ana", "app", 6, 5, et.Json{})
		}},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			_, calls, works := withFakes(t, nil, et.Item{Ok: true, Result: et.Json{"count": 1}}, nil)
			result := c.run()
			if len(*works) != 0 {
				t.Fatalf("no debe publicar event.Work: %v", *works)
			}
			if len(*calls) != 1 || (*calls)[0].method != c.method {
				t.Fatalf("llamada = %v, quiere %s", *calls, c.method)
			}
			serviceId := (*calls)[0].args.Str("service_id")
			if serviceId == "" || result.Str("service_id") != serviceId {
				t.Fatalf("service_id = %q, enviado %q", result.Str("service_id"), serviceId)
			}
			if !result.Bool("ok") || result.Int("count") != 1 {
				t.Fatalf("resultado = %s", result.ToString())
			}
			if _, ok := result["error"]; ok {
				t.Fatalf("no debe tener error: %s", result.ToString())
			}
		})
	}
}

func TestSendRpcEnvPackage(t *testing.T) {
	t.Setenv("SEND_RPC", "Otro-send.Services")
	t.Setenv("OTP_RPC", "Otro-otp.Services")
	_, calls, _ := withFakes(t, nil, et.Item{Ok: true}, nil)
	SendSms("p1", "", []string{"3000000000"}, "hola", nil, TpTransactional, et.Json{})
	SendOtpEmail("p1", "a@b.co", "Ana", "app", 6, 5, et.Json{})
	if (*calls)[0].method != "Otro-send.Services.SendSMS" || (*calls)[1].method != "Otro-otp.Services.GenerateOTPEmail" {
		t.Fatalf("metodos = %v", *calls)
	}
	if (*calls)[0].args.Str("service_id") == "" {
		t.Fatal("service_id vacio debe generarse")
	}
}

func TestSendRpcFallbackWhenNotFound(t *testing.T) {
	_, calls, works := withFakes(t, fmt.Errorf(jrpc.ERR_METHOD_NOT_FOUND, "Service-send.Services.SendSMS"), et.Item{}, nil)
	result := SendSms("p1", "s1", []string{"3000000000"}, "hola", nil, TpTransactional, et.Json{"id": "u1"})
	if len(*calls) != 0 {
		t.Fatalf("no debe llamar el RPC: %v", *calls)
	}
	if len(*works) != 1 || (*works)[0].method != "send/sms" {
		t.Fatalf("event.Work = %v", *works)
	}
	data := (*works)[0].args
	if data.Str("content") != "hola" || data.Str("type") != "transactional" || data.Str("project_id") != "p1" {
		t.Fatalf("data del evento = %s", data.ToString())
	}
	if result.Str("service_id") != "s1" || result.Str("event") != "send/sms" {
		t.Fatalf("resultado = %s", result.ToString())
	}

	_, _, works = withFakes(t, errors.New(jrpc.ERR_PACKAGE_NOT_FOUND), et.Item{}, nil)
	SendOtpWhatsapp("p1", "3000000000", "Ana", "app", 6, 5, et.Json{})
	if len(*works) != 1 || (*works)[0].method != "generate/otp/whatsapp" || (*works)[0].args.Str("kind") != "whatsapp" {
		t.Fatalf("event.Work otp = %v", *works)
	}
}

func TestSendRpcCallErrorNoResend(t *testing.T) {
	_, calls, works := withFakes(t, nil, et.Item{}, errors.New("proveedor caido"))
	result := SendEmail("p1", "s1", []et.Json{{"email": "a@b.co"}}, "asunto", "<p/>", nil, TpTransactional, et.Json{})
	if len(*calls) != 1 {
		t.Fatalf("llamadas = %v", *calls)
	}
	if len(*works) != 0 {
		t.Fatalf("no debe reenviar por event.Work: %v", *works)
	}
	if result.Bool("ok") || result.Str("error") != "proveedor caido" || result.Str("service_id") != "s1" {
		t.Fatalf("resultado = %s", result.ToString())
	}
}

func TestRpcResultOkFalseWithoutError(t *testing.T) {
	result := rpcResult(et.Item{Ok: false, Result: et.Json{"ok": true, "count": 0}}, nil, "s1")
	if result.Bool("ok") || result.Str("service_id") != "s1" {
		t.Fatalf("resultado = %s", result.ToString())
	}
	if _, ok := result["error"]; ok {
		t.Fatalf("sin error Go no debe tener error: %s", result.ToString())
	}
}

// Services de prueba: replica la firma de los RPC de send (require et.Json, response *et.Item) error
type Services struct{ last et.Json }

func (s *Services) SendSMS(require et.Json, response *et.Item) error {
	s.last = require
	if require.Str("content") == "falla" {
		return errors.New("Error: proveedor rechazo el mensaje")
	}
	*response = et.Item{Ok: true, Result: et.Json{"ok": true, "count": 1, "result": []et.Json{{"to": require.ArrayStr("contact_numbers")[0]}}}}
	return nil
}

func TestSendRpcOverTCP(t *testing.T) {
	server := rpc.NewServer()
	services := &Services{}
	if err := server.Register(services); err != nil {
		t.Fatal(err)
	}
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	defer listener.Close()
	go server.Accept(listener)
	port := listener.Addr().(*net.TCPAddr).Port

	oldSolve, oldCall, oldWork := rpcSolve, rpcCall, workPublish
	t.Cleanup(func() { rpcSolve, rpcCall, workPublish = oldSolve, oldCall, oldWork })
	// jrpc.call registra metricas en Redis (cache.Incr), que la prueba no tiene: el transporte es el
	// mismo net/rpc + gob, sin la metrica
	rpcCall = func(solver *jrpc.Solver, args et.Json) (et.Item, error) {
		var item et.Item
		client, err := rpc.Dial("tcp", fmt.Sprintf("%s:%d", solver.Host, solver.Port))
		if err != nil {
			return item, err
		}
		defer client.Close()
		err = client.Call(solver.Method, args, &item)
		return item, err
	}
	// Como lo registra jrpc.Mount: Method = <Struct>.<Metodo>, Host/Port del paquete
	rpcSolve = func(method string) (*jrpc.Solver, error) {
		return &jrpc.Solver{Method: "Services.SendSMS", Host: "127.0.0.1", Port: port}, nil
	}
	workPublish = func(event string, data et.Json) et.Json {
		t.Fatalf("no debe publicar event.Work")
		return nil
	}

	result := SendSms("p1", "s1", []string{"3000000000"}, "hola", []et.Json{{"code": "123"}}, TpTransactional, et.Json{"id": "u1"})
	if !result.Bool("ok") || result.Int("count") != 1 || result.Str("service_id") != "s1" {
		t.Fatalf("resultado = %s", result.ToString())
	}
	if services.last.Str("type") != "transactional" || services.last.Json("user").Str("id") != "u1" {
		t.Fatalf("args recibidos = %s", services.last.ToString())
	}

	result = SendSms("p1", "s2", []string{"3000000000"}, "falla", nil, TpTransactional, et.Json{})
	if result.Bool("ok") || result.Str("error") != "Error: proveedor rechazo el mensaje" || result.Str("service_id") != "s2" {
		t.Fatalf("resultado con error = %s", result.ToString())
	}
}

func TestSendRpcFallbackRealSolverWithoutRedis(t *testing.T) {
	// Sin Redis (cache.Load) jrpc.GetSolver devuelve error: no hubo llamada, se publica el evento
	t.Setenv("PIPE_HOST", "")
	oldWork := workPublish
	t.Cleanup(func() { workPublish = oldWork })
	works := []string{}
	workPublish = func(event string, data et.Json) et.Json {
		works = append(works, event)
		return et.Json{"_id": data.Str("service_id"), "event": event, "data": data}
	}

	result := SendEmailToByTemplate("p1", "s1", "Ana", "a@b.co", "asunto", "t1", nil, TpTransactional, et.Json{})
	if len(works) != 1 || works[0] != "send/emailtobytemplate" || result.Str("service_id") != "s1" {
		t.Fatalf("event.Work = %v resultado = %s", works, result.ToString())
	}
}
