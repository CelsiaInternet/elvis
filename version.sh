#!/bin/bash

set -e                                                  # Detener la ejecución en caso de error

MAYOR=false
MINOR=false
INDEX=2
VERSION=$(git describe --tags --abbrev=0)               # Valor para reemplazar $VERSION obtenido de Git
README_FILE="README.md"                                 # Nombre del archivo a modificar
HELP=false                                              # Bandera para mostrar la ayuda

# Parsear opciones
while [[ "$#" -gt 0 ]]; do
    case $1 in
        --h | --help) HELP=true ;;                      # Mostrar la ayuda
        --m | --major) MAYOR=true ;;                    # Activar la bandera si se proporciona --major
        --n | --minor) MINOR=true ;;                    # Activar la bandera si se proporciona --minor
        *) echo "Opción desconocida: $1"; exit 1 ;;
    esac
    shift
done

# Mostrar las opciones elegidas
echo "Opciones elegidas:"
[[ "$MAYOR" == true ]] && echo " - Major: Activado"
[[ "$MINOR" == true ]] && echo " - Minor: Activado"

# Obtiene la última etiqueta
latest_tag=$(git describe --tags --abbrev=0 2>/dev/null)

# Si no hay etiquetas, usa la versión predeterminada v1.0.0
if [ -z "$latest_tag" ]; then
  new_version="v1.0.0"
elif [ "$HELP" == true ]; then
  echo "Uso: ./version.sh [opciones]"
  echo "Opciones:"
  echo "  --h, --help     Mostrar esta ayuda"
  echo "  --m, --major    Incrementar la versión mayor"
  echo "  --n, --minor    Incrementar la versión menor"
  exit 0
else
  # Divide la etiqueta en componentes usando el punto como delimitador
  IFS='.' read -r -a version_parts <<< "${latest_tag#v}"

  if [ "$MAYOR" == true ]; then
    # Si se proporciona la opción --major, incrementa el valor de la posición 0
    version_parts[0]=$((version_parts[0] + 1))
    version_parts[1]=0
    version_parts[2]=0
  
    new_version="v${version_parts[0]}.${version_parts[1]}.${version_parts[2]}"
  elif [ "$MINOR" == true ]; then
    # Si se proporciona la opción --minor, incrementa el valor de la posición 1        
    version_parts[1]=$((version_parts[1] + 1))
    version_parts[2]=0

    new_version="v${version_parts[0]}.${version_parts[1]}.${version_parts[2]}"
  else
    # Incrementa el valor de la posición 2
    version_parts[2]=$((version_parts[2] + 1))

    # Reconstruye la nueva versión (X.Y.Z) y prepende la 'v' al principio
    new_version="v${version_parts[0]}.${version_parts[1]}.${version_parts[2]}"
  fi  
fi

sed -i '' "s/$VERSION/$new_version/g" "$README_FILE"

# Muestra la nueva versión
echo "La nueva versión es: $new_version"

git add .
git commit -m 'Update'
git push -u origin main
git tag "$new_version"
git push -u origin --tags