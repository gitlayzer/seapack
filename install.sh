#!/usr/bin/env bash

# Repository configuration
REPO_OWNER="gitlayzer"
REPO_NAME="seapack"
REPO_URL="https://github.com/${REPO_OWNER}/${REPO_NAME}"

help_text="Options

   -V, --verbose
   Enable verbose output for the installer

   -f, -y, --force, --yes
   Skip the confirmation prompt during installation

   -p, --platform
   Override the platform identified by the installer

   -b, --bin-dir
   Override the bin installation directory

   -a, --arch
   Override the architecture identified by the installer

   -B, --base-url
   Override the base URL used for downloading releases

   -r, --remove
   Uninstall seapack

   -h, --help
   Get some help

"

set -eu
printf '\n'

BOLD="$(tput bold 2>/dev/null || printf '')"
GREY="$(tput setaf 0 2>/dev/null || printf '')"
UNDERLINE="$(tput smul 2>/dev/null || printf '')"
RED="$(tput setaf 1 2>/dev/null || printf '')"
GREEN="$(tput setaf 2 2>/dev/null || printf '')"
YELLOW="$(tput setaf 3 2>/dev/null || printf '')"
BLUE="$(tput setaf 4 2>/dev/null || printf '')"
MAGENTA="$(tput setaf 5 2>/dev/null || printf '')"
NO_COLOR="$(tput sgr0 2>/dev/null || printf '')"


SUPPORTED_TARGETS="linux-amd64 linux-arm64"

CURL_RETRY_OPTS="--retry 3 --retry-all-errors --retry-delay 2"

info() {
  printf '%s\n' "${BOLD}${GREY}>${NO_COLOR} $*"
}

debug() {
  if [[ -n "${VERBOSE}" ]]; then
    printf '%s\n' "${BOLD}${GREY}>${NO_COLOR} $*"
  fi
}

warn() {
  printf '%s\n' "${YELLOW}! $*${NO_COLOR}"
}

error() {
  printf '%s\n' "${RED}x $*${NO_COLOR}" >&2
}

completed() {
  printf '%s\n' "${GREEN}✓${NO_COLOR} $*"
}

has() {
  command -v "$1" 1>/dev/null 2>&1
}

get_tmpfile() {
  local suffix
  suffix="$1"
  if has mktemp; then
    printf "%s%s.%s.%s" "$(mktemp)" "-seapack" "$(date +%s)" "${suffix}"
  else
    printf "/tmp/seapack.%s" "${suffix}"
  fi
}

# Test if a location is writeable by trying to write to it.
test_writeable() {
  local path
  path="${1:-}/test.txt"
  if touch "${path}" 2>/dev/null; then
    rm -f "${path}"
    return 0
  else
    return 1
  fi
}

download() {
 file="$1"
  url="$2"
  touch "$file"

  if has curl; then
    cmd="curl --fail --silent --location ${CURL_RETRY_OPTS} --output $file $url"
  elif has wget; then
    cmd="wget --quiet --output-document=$file $url"
  elif has fetch; then
    cmd="fetch --quiet --output=$file $url"
  else
    error "No HTTP download program (curl, wget, fetch) found, exiting…"
    return 1
  fi

  $cmd && return 0 || rc=$?

  error "Command failed (exit code $rc): ${BLUE}${cmd}${NO_COLOR}"
  printf "\n" >&2
  info "This is likely due to seapack not yet supporting your configuration."
  info "If you would like to see a build for your configuration,"
  info "please create an issue requesting a build for ${MAGENTA}${TARGET}${NO_COLOR}:"
  info "${BOLD}${UNDERLINE}${REPO_URL}/issues/new/${NO_COLOR}"
  return $rc
}

unpack() {
  local archive=$1
  local bin_dir=$2
  local sudo=${3-}

  case "$archive" in
    *.tar.gz)
      flags=$(test -n)
      ${sudo} tar "${flags}" -xzf "${archive}" -C "${bin_dir}"
      return 0
      ;;
    *.zip)
      flags=$(test -z)
      UNZIP="${flags}" ${sudo} unzip "${archive}" -d "${bin_dir}"
      return 0
      ;;
  esac

  error "Unknown package extension."
  printf "\n"
  info "This almost certainly results from a bug in this script--please file a"
  info "bug report at ${REPO_URL}/issues"
  return 1
}

elevate_priv() {
  if ! has sudo; then
    error 'Could not find the command "sudo", needed to get permissions for install.'
    info "Run this script as root, install sudo, or choose a writable --bin-dir."
    exit 1
  fi
  if ! sudo -v; then
    error "Superuser not granted, aborting installation"
    exit 1
  fi
}

install() {
  local msg
  local sudo
  local archive
  local ext="$1"

  if test_writeable "${BIN_DIR}"; then
    sudo=""
    msg="Installing seapack, please wait…"
  else
    warn "Escalated permissions are required to install to ${BIN_DIR}"
    elevate_priv
    sudo="sudo"
    msg="Installing seapack as root, please wait…"
  fi
  info "$msg"

  archive=$(get_tmpfile "$ext")

  # download to the temp file
  download "${archive}" "${URL}"

  # unpack the temp file to the bin dir, using sudo if required
  unpack "${archive}" "${BIN_DIR}" "${sudo}"

  # remove tempfile
  # rm "${archive}"
}

# Simplify platform detection
detect_platform() {
  local platform
  platform="$(uname -s | tr '[:upper:]' '[:lower:]')"

  case "${platform}" in
    linux) platform="linux" ;;
  esac

  printf '%s' "${platform}"
}

detect_arch() {
  local arch
  arch="$(uname -m | tr '[:upper:]' '[:lower:]')"

  case "${arch}" in
    amd64|x86_64) printf 'amd64' ;;
    arm64|aarch64) printf 'arm64' ;;
    *) printf '%s' "${arch}" ;;
  esac
}

detect_target() {
  local arch="$1"
  local platform="$2"

  printf '%s-%s' "${platform}" "${arch}"
}


confirm() {
  if [ -t 0 ]; then
    if [ -z "${FORCE-}" ]; then
      printf "%s " "${MAGENTA}?${NO_COLOR} $* ${BOLD}[y/N]${NO_COLOR}"
      set +e
      read -r yn </dev/tty
      rc=$?
      set -e
      if [ $rc -ne 0 ]; then
        error "Error reading from prompt (please re-run with the '--yes' option)"
        exit 1
      fi
      if [ "$yn" != "y" ] && [ "$yn" != "yes" ]; then
        error 'Aborting (please answer "yes" to continue)'
        exit 1
      fi
    fi
  fi
}

check_bin_dir() {
  local bin_dir="$1"

  if [ ! -d "$BIN_DIR" ]; then
    error "Installation location $BIN_DIR does not appear to be a directory"
    info "Make sure the location exists and is a directory, then try again."
    exit 1
  fi

  # https://stackoverflow.com/a/11655875
  local good
  good=$(
    IFS=:
    for path in $PATH; do
      if [ "${path}" = "${bin_dir}" ]; then
        printf 1
        break
      fi
    done
  )

  if [ "${good}" != "1" ]; then
    warn "Bin directory ${bin_dir} is not in your \$PATH"
  fi
}

is_build_available() {
  local arch="$1"
  local platform="$2"
  local target="$3"

  local good

  good=$(
    IFS=" "
    for t in $SUPPORTED_TARGETS; do
      if [ "${t}" = "${target}" ]; then
        printf 1
        break
      fi
    done
  )

  if [ "${good}" != "1" ]; then
    error "${arch} builds for ${platform} are not yet available for seapack"
    printf "\n" >&2
    info "If you would like to see a build for your configuration,"
    info "please create an issue requesting a build for ${MAGENTA}${target}${NO_COLOR}:"
    info "${BOLD}${UNDERLINE}${REPO_URL}/issues/new/${NO_COLOR}"
    printf "\n"
    exit 1
  fi
}

UNINSTALL=0
HELP=0

# defaults
if [ -z "${SEAPACK_VERSION-}" ]; then
  SEAPACK_VERSION=$(curl --fail --silent --show-error ${CURL_RETRY_OPTS} ${GITHUB_TOKEN:+-H "Authorization: Bearer $GITHUB_TOKEN"} "https://api.github.com/repos/${REPO_OWNER}/${REPO_NAME}/releases/latest" | grep -o '"tag_name": "v.*"' | cut -d'"' -f4 | cut -c2-)

  if [ -z "${SEAPACK_VERSION}" ]; then
    error "Failed to detect latest version from GitHub API"
    info "You can set SEAPACK_VERSION=x.y.z to specify a version manually"
    exit 1
  fi
fi

if [ -n "${SEAPACK_PLATFORM-}" ]; then
  PLATFORM="${SEAPACK_PLATFORM}"
else
  PLATFORM="$(detect_platform)"
fi

if [ -n "${SEAPACK_BIN_DIR-}" ]; then
  BIN_DIR="${SEAPACK_BIN_DIR}"
else
  BIN_DIR=/usr/local/bin
fi

if [ -n "${SEAPACK_ARCH-}" ]; then
  ARCH="${SEAPACK_ARCH}"
else
  ARCH="$(detect_arch)"
fi

if [ -n "${SEAPACK_BASE_URL-}" ]; then
  BASE_URL="${SEAPACK_BASE_URL}"
else
  BASE_URL="${REPO_URL}/releases"
fi

while [ "$#" -gt 0 ]; do
  case "$1" in
  -p | --platform)
    PLATFORM="$2"
    shift 2
    ;;
  -b | --bin-dir)
    BIN_DIR="$2"
    shift 2
    ;;
  -a | --arch)
    ARCH="$2"
    shift 2
    ;;
  -B | --base-url)
    BASE_URL="$2"
    shift 2
    ;;

  -V | --verbose)
    VERBOSE=1
    shift 1
    ;;
  -f | -y | --force | --yes)
    FORCE=1
    shift 1
    ;;
  -r | --remove | --uninstall)
    UNINSTALL=1
    shift 1
    ;;
  -h | --help)
    HELP=1
    shift 1
    ;;
  -p=* | --platform=*)
    PLATFORM="${1#*=}"
    shift 1
    ;;
  -b=* | --bin-dir=*)
    BIN_DIR="${1#*=}"
    shift 1
    ;;
  -a=* | --arch=*)
    ARCH="${1#*=}"
    shift 1
    ;;
  -B=* | --base-url=*)
    BASE_URL="${1#*=}"
    shift 1
    ;;
  -V=* | --verbose=*)
    VERBOSE="${1#*=}"
    shift 1
    ;;
  -f=* | -y=* | --force=* | --yes=*)
    FORCE="${1#*=}"
    shift 1
    ;;

  *)
    error "Unknown option: $1"
    exit 1
    ;;
  esac
done

if [ -n "${VERBOSE-}" ]; then
  VERBOSE=v
else
  VERBOSE=
fi

if [ $UNINSTALL == 1 ]; then
  confirm "Are you sure you want to uninstall seapack?"

  msg=""
  sudo=""

  if test_writeable "$(dirname "$(which ${REPO_NAME})")"; then
    sudo=""
    msg="Removing seapack, please wait…"
  else
    warn "Escalated permissions are required to install to ${BIN_DIR}"
    elevate_priv
    sudo="sudo"
    msg="Removing seapack as root, please wait…"
  fi

  info "$msg"
  ${sudo} rm -f "$(which ${REPO_NAME})"
  ${sudo} rm -rf /tmp/${REPO_NAME}

  info "Removed seapack"
  exit 0
fi

if [ $HELP == 1 ]; then
    echo "${help_text}"
    exit 0
fi
TARGET="$(detect_target "${ARCH}" "${PLATFORM}")"

is_build_available "${ARCH}" "${PLATFORM}" "${TARGET}"

print_configuration () {
  if [[ -n "${VERBOSE}" ]]; then
    printf "  %s\n" "${UNDERLINE}Configuration${NO_COLOR}"
    debug "${BOLD}Bin directory${NO_COLOR}: ${GREEN}${BIN_DIR}${NO_COLOR}"
    debug "${BOLD}Platform${NO_COLOR}:      ${GREEN}${PLATFORM}${NO_COLOR}"
    debug "${BOLD}Arch${NO_COLOR}:          ${GREEN}${ARCH}${NO_COLOR}"
    debug "${BOLD}Version${NO_COLOR}:       ${GREEN}${SEAPACK_VERSION}${NO_COLOR}"
    printf '\n'
  fi
}

print_configuration

EXT=tar.gz
URL="${BASE_URL}/download/v${SEAPACK_VERSION}/${REPO_NAME}-${TARGET}.${EXT}"
debug "Tarball URL: ${UNDERLINE}${BLUE}${URL}${NO_COLOR}"
confirm "Install seapack ${GREEN}${SEAPACK_VERSION}${NO_COLOR} to ${BOLD}${GREEN}${BIN_DIR}${NO_COLOR}?"
check_bin_dir "${BIN_DIR}"

install "${EXT}"

printf "$GREEN"
cat <<'EOF'

   ____       _ _                  _
  |  _ \ __ _(_) |_ __   __ _  ___| | __
  | |_) / _` | | | '_ \ / _` |/ __| |/ /
  |  _ < (_| | | | |_) | (_| | (__|   <
  |_| \_\__,_|_|_| .__/ \__,_|\___|_|\_\
                  |_|

  SeaPack is now installed!
  Run 'seapack --help' to get started

EOF
printf "$NO_COLOR"
