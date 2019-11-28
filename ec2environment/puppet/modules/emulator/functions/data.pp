function emulator::data(
  String $key,
  Optional[String] $value=undef
) {
  if $value {
    info("Caching ${key} = '${value}' for future use")
    choria::data($key, $value, emulator::ds())
  } else {
    choria::data($key, emulator::ds())
  }
}