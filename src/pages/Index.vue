<template>
  <q-page padding>
    Bluetooth devices
    <q-list bordered separator>
      <q-item v-for="(device,index) in devices" :key="index" @click.native="connectDevice(device)" clickable>
        <q-item-section>Name: {{ device.title }} - ID: {{device.id}}</q-item-section>
      </q-item>
    </q-list>
    <q-btn :label="isScaning?'Scanning':'Scan'" @click.native="scanDevices" :disabled="isScaning"></q-btn>
    {{deviceCharac}}
    <br/>
  </q-page>
</template>

<script>
export default {
  data(){
    return {
      devices:[],
      subscription: null,
      isScaning: false,
      deviceCharac: ""
          }
  },
  methods:{
    stopScaning(){
      this.isScaning= false
      psgo.call("web.stop.bluetooth.devices").then( result =>{
        for(var i=0;i<result.length;i++){          
          this.devices.push({title:result[i][0], id:result[i][1]})
        }
      })

    }, 
    scanDevices(){
      let _this = this
      this.isScaning = true
      this.devices = []
      psgo.pub("web.scan.bluetooth.devices", null)
      setTimeout(function(){
        _this.stopScaning()
      },3000)
    },
    connectDevice(dev){
      this.$router.go({name: "char"})
      /*console.log("connect")
      if(dev.title!="Error"){
        console.log("connect2")
        psgo.call("web.connect.bluetooth.devices", dev.id).then( result =>{
          console.log("connected")
          console.log("result")
          this.deviceCharac = result
        })
      }*/
    }
  }
}
</script>

<style>

</style>
