plan emulator::start (
  String $servers,
  Boolean $tls=false,
  Optional[String] $credentials=undef,
  Integer $monitor=8080,
  Optional[String] $region=undef,
  Optional[String] $iname=undef,
  Integer $instances=1,
  Integer $agents=1,
  Integer $collectives=1,
  Array[String] $nodes=[]
) {
  if !$iname and $region {
    $_name = $region
  } elsif $iname {
    $_name = $iname
  } else {
    $_name = "emulator"
  }
 
  if $credentials {
    $_cred_properties = {
      "credentials" => base64(encode, file($credentials))
    }
  } else {
    $_cred_properties = {}
  }

  $_base_properties = {
    "name" => $_name,
    "servers" => $servers,
    "tls" => $tls,
    "monitor" => $monitor,
    "instances" => $instances,
    "agents" => $agents,
    "collectives" => $collectives 
  }

  if $nodes.length > 0 {
    $_nodes = $nodes
  } elsif $region {
    $_nodes = emulator::regional_nodes($region)
  } else {
    $_nodes = emulator::all_nodes()
  }

  $results = choria::task("mcollective",
    "nodes" => $_nodes,
    "action" => "emulator.start",
    "silent" => true,
    "post" => ["summarize"],
    "properties" => $_base_properties + $_cred_properties
  ) 

  $results.each |$result| {
    if $result.ok {
      info(sprintf("%s: %s: started: %s", $result["sender"], $result["statusmsg"], $result["data"]["status"]))
    } else {
      error(sprintf("%s: %s", $result["sender"], $result["statusmsg"]))
    }
  }

  undef
}

