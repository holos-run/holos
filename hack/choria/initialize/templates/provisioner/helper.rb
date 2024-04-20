#!/usr/bin/ruby

require "json"
require "yaml"
require "base64"
require "net/http"
require "openssl"

def parse_input
  input = STDIN.read

  begin
    File.open("/tmp/request.json", "w") {|f| f.write(input)}
  rescue Exception
  end

  request = JSON.parse(input)
  request["inventory"] = JSON.parse(request["inventory"])

  request
end

def validate!(request, reply)
  if request["identity"] && request["identity"].length == 0
    reply["msg"] = "No identity received in request"
    reply["defer"] = true
    return false
  end

  unless request["ed25519_pubkey"]
    reply["msg"] = "No ed15519 public key received"
    reply["defer"] = true
    return false
  end

  unless request["ed25519_pubkey"]
    reply["msg"] = "No ed15519 directory received"
    reply["defer"] = true
    return false
  end

  if request["ed25519_pubkey"]["directory"].length == 0
    reply["msg"] = "No ed15519 directory received"
    reply["defer"] = true
    return false
  end

  true
end

def publish_reply(reply)
  begin
    File.open("/tmp/reply.json", "w") {|f| f.write(reply.to_json)}
  rescue Exception
  end

  puts reply.to_json
end

def publish_reply!(reply)
  publish_reply(reply)
  exit
end

def set_config!(request, reply)
  # stub data the helper will fetch from the saas
  customers = {
        "one" => {
            :brokers => "nats://managed.example.net:9222", # whoever is the leader for this site
            :site => "customer_one",
            :source => {
                :host => "nats://cust_one:s3cret@saas-nats.choria.local",
            }
        }
    }


  customer = request["jwt"]["extensions"]["customer"]
  brokers = customers[customer][:brokers]
  source = customers[customer][:source]

  reply["configuration"].merge!(
    "identity" => request["identity"],
    "loglevel" => "warn",
    "plugin.choria.server.provision" => "false",
    "plugin.choria.middleware_hosts" => brokers,
    "plugin.security.issuer.names" => "choria",
    "plugin.security.issuer.choria.public" => "ISSUER",
    "plugin.security.provider" => "choria",
    "plugin.security.choria.token_file" => File.join(request["ed25519_pubkey"]["directory"], "server.jwt"),
    "plugin.security.choria.seed_file" => File.join(request["ed25519_pubkey"]["directory"], "server.seed"),
    "machine_room.role" => "leader",
    "machine_room.site" => customers[customer][:site],
    "machine_room.source.host" => source[:host],
  )

  reply["server_claims"].merge!(
    "exp" => 5*60*60*24*365,
    "pub_subjects" => [">"],
    "permissions" => {
      "streams" => true,
      "submission" => true,
      "service_host" => true,
    }
  )
end

reply = {
  "defer" => false,
  "msg" => "",
  "certificate" => "",
  "ca" => "",
  "configuration" => {},
  "server_claims" => {}
}

begin
  request = parse_input

  reply["msg"] = "Validating"
  unless validate!(request, reply)
    publish_reply!(reply)
  end

  set_config!(request, reply)

  reply["msg"] = "Done"
  publish_reply!(reply)
rescue SystemExit
rescue Exception
  reply["msg"] = "Unexpected failure during provisioning: %s: %s" % [$!.class, $!.to_s]
  reply["defer"] = true
  publish_reply!(reply)
end
