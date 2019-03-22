<template>
  <q-page padding>
    Bluetooth devices
    <q-list bordered separator>
      <q-item v-for="(device,index) in devices" :key="index" @click.native="connectDevice(device)" clickable>
        <q-item-section>Name: {{ device.title }} - ID: {{device.id}}</q-item-section>
      </q-item>
    </q-list>
    <q-btn :label="isScaning?'Scanning':'Scan'" @click.native="scanDevices" :disabled="isScaning"></q-btn>
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
      if(dev.title!="Error"){
        this.$router.push({name:"char", params:{title: dev.title, id: dev.id}})
      }
    }
  }
}
</script>

<style>

</style>
