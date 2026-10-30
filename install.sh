#!/usr/bin/env bash
#
# install.sh - Instala cjk-press en Linux desde los releases de GitHub.
#
# Uso:
#   curl -fsSL https://raw.githubusercontent.com/atetov/cjk-press/main/install.sh | bash
#   INSTALL_DIR=$HOME/.local/bin VERSION=v0.1.0 bash install.sh
#
set -euo pipefail

REPO="${REPO:-atetov/cjk-press}"
BIN_NAME="cjk-press"
INSTALL_DIR="${INSTALL_DIR:-/usr/local/bin}"
VERSION="${VERSION:-latest}"

# Detectar arquitectura.
ARCH="$(uname -m)"
case "${ARCH}" in
  x86_64)  ARCH="amd64" ;;
  aarch64) ARCH="arm64" ;;
  armv*)   ARCH="arm64" ;;
  *)       echo "Arquitectura no soportada: ${ARCH}" >&2; exit 1 ;;
esac

OS="linux"
ARTIFACT="${BIN_NAME}_${VERSION#v}_${OS}_${ARCH}.tar.gz"
# Para "latest" consultamos la API de GitHub para obtener el tag real.
if [ "${VERSION}" = "latest" ]; then
  TAG="$(curl -fsSL "https://api.github.com/repos/${REPO}/releases/latest" \
        | grep '"tag_name":' | sed -E 's/.*"([^"]+)".*/\1/')"
  if [ -z "${TAG}" ]; then
    echo "No se pudo obtener la versión latest. Especifica VERSION=v0.1.0" >&2
    exit 1
  fi
  VERSION="${TAG}"
fi
ARTIFACT="${BIN_NAME}_${VERSION#v}_${OS}_${ARCH}.tar.gz"
URL="https://github.com/${REPO}/releases/download/${VERSION}/${ARTIFACT}"

echo "==> cjk-press installer"
echo "    Repo:     ${REPO}"
echo "    Versión:  ${VERSION}"
echo "    Arch:     ${OS}/${ARCH}"
echo "    Destino:  ${INSTALL_DIR}"

TMPDIR="$(mktemp -d)"
trap 'rm -rf "${TMPDIR}"' EXIT

echo "==> Descargando ${URL}..."
curl -fsSL -o "${TMPDIR}/${ARTIFACT}" "${URL}"

echo "==> Extrayendo..."
tar -xzf "${TMPDIR}/${ARTIFACT}" -C "${TMPDIR}"

# El binario dentro del tarball se llama cjk-press.
if [ ! -f "${TMPDIR}/${BIN_NAME}" ]; then
  echo "No se encontró el binario ${BIN_NAME} en el paquete." >&2
  exit 1
fi

echo "==> Instalando en ${INSTALL_DIR}/${BIN_NAME}..."
if [ ! -w "${INSTALL_DIR}" ]; then
  echo "Se necesitan permisos de root para escribir en ${INSTALL_DIR}."
  sudo install -m 0755 "${TMPDIR}/${BIN_NAME}" "${INSTALL_DIR}/${BIN_NAME}"
else
  install -m 0755 "${TMPDIR}/${BIN_NAME}" "${INSTALL_DIR}/${BIN_NAME}"
fi

echo "==> ¡Listo! Prueba con:"
echo "    ${BIN_NAME} --version"
echo "    ${BIN_NAME} --help"