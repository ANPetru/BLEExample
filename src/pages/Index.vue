<template>
  <q-page padding>
    <div v-if="devices!=null">
      
    <q-toolbar class="bg-secondary text-white shadow-2 " >
      <q-toolbar-title>Bluetooth devices</q-toolbar-title>
    </q-toolbar>
    <q-list bordered separator>
      <q-item v-for="(device,index) in devices" :key="index" @click.native="connectDevice(device)" clickable v-bind:class="{listItem1:index%2==0, listItem2:index%2==1}">
        <q-item-section>Name: {{ device.title }} - ID: {{device.id}}</q-item-section>

      </q-item>
    </q-list>
    </div>
      <q-btn v-if="!isScaning" label="Scan" @click.native="scanDevices" :disabled="isScaning" color="primary" style="position:fixed; bottom:10px; left:40%; width:20%;" ></q-btn>
  </q-page>
</template>

<script>
export default {
  data(){
    return {
      devices:null,
      subscription: null,
      isScaning: false,
          }
  },
  methods:{
    stopScaning(){
      this.isScaning= false
      psgo.call("web.stop.bluetooth.devices").then( result =>{
        if(result.length>0){
          this.devices=[]
          for(var i=0;i<result.length;i++){          
            this.devices.push({title:result[i][0], id:result[i][1], rssi: result[i][2]})
            console.log(result[i][2])
          }
        }
      })

    }, 
    scanDevices(){
      this.showLoading()
      let _this = this
      this.isScaning = true
      this.devices = null
      psgo.pub("web.scan.bluetooth.devices", null)
      setTimeout(function(){
        _this.stopScaning()
      },3000)
    },
    connectDevice(dev){
      if(dev.title!="Error"){
        this.$router.push({name:"char", params:{title: dev.title, id: dev.id}})
      }
    },
    showLoading(){ 
      this.$q.loading.show({
        message: "Scanning bluetooth devices"
      })
      this.timer = setTimeout(() => {
      this.$q.loading.hide()
      this.timer = void 0}, 3000)
    }
  }
}
</script>

<style>


.listItem1{
  padding-left: 30px;
  background-color: rgb(240, 255, 250);
}
.listItem2{
  padding-left: 30px;
  background-color: rgb(241, 255, 240);
}
</style>
