plan emulator::nats::start (
  Boolean $leafnode=false,
  Optional[String] $servers=undef,
  Optional[String] $credentials=undef,
  Integer $monitor=8222,
  Integer $clients=4222,
  Array[String] $nodes
) {
  if $leafnode {
    if !$servers {
        fail("leaf nodes need servers specified")
    }

    info("Starting leafnode nats-servers")

    if $credentials and $credentials != "" {
      $_cred_options = {
        "credentials" => base64(encode, file($credentials)),
      }
    } else {
      $_cred_options = {}
    }

    $_leaf_options = $_cred_options + {
        "leafnode_servers" => $servers
    }
  } else {
    info("Starting standalone nats-server")
    $_leaf_options = {}
  }
  
  $_common_options = {
    "monitor_port" => $monitor,
    "port" => $clients
  }

  $results = choria::task("mcollective",
    "nodes" => $nodes,
    "action" => "emulator.start_nats",
    "silent" => true,
    "properties" => $_common_options + $_leaf_options
  ) 

  $results.each |$result| {
    if $result.ok {
      info(sprintf("%s: %s: started: %s", $result["sender"], $result["statusmsg"], $result["data"]["running"]))
    } else {
      error(sprintf("%s: %s", $result["sender"], $result["statusmsg"]))
    }
  }

  undef
}