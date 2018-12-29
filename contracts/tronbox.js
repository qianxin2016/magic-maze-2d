module.exports = {
  networks: {
    development: {
      privateKey: 'bb1a42f78ee87014921ea24ce572947f243077dde4c68f4f893f80a66eb1f03f',
      consume_user_resource_percent: 30,
      fee_limit: 100000000,

      fullNode: "https://api.shasta.trongrid.io",
      solidityNode: "https://api.shasta.trongrid.io",
      eventServer: "https://api.shasta.trongrid.io",

      network_id: "*"
    }
  }
}
