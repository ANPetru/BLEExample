<template>
    <q-page padding>
        <div v-if="deviceCharac!=''">
            <q-toolbar class="bg-secondary text-white shadow-2">
            <q-toolbar-title>Connected To: {{connectedDevice.title}}</q-toolbar-title>
            Signal:{{signal}}
        </q-toolbar>
            <q-list bordered separator>
                <q-item v-for="(charac,index) in deviceCharac" :key="index" v-bind:class="{listItem1:index%2==0, listItem2:index%2==1}">
                    <q-item-section class="col">Service: {{ charac.Service }} - {{charac.Characteristic}}
                    </q-item-section>
                    <q-item-section side>
                        <q-icon class="col" v-if="canReadCharac(charac)" name="visibility" @click.native="readCharac(charac)"></q-icon>
                    </q-item-section>
                    <q-item-section side>
                        <q-icon class="col" name="create" @click.native="prepareMessage(charac)"></q-icon>
                    </q-item-section>
                </q-item>
            </q-list>
        </div>
    </q-page>
</template>

<script>
export default {
    data(){
        return {
            connectedDevice: {title:this.$route.params.title, id:this.$route.params.id},
            deviceCharac: '',
            signal:0
        }
    }, 
    methods:{
        connectToDevice(){
            this.$q.loading.show({
                message: "Conntecting to " + this.connectedDevice.title
            })
            psgo.call("web.connect.bluetooth.devices", this.connectedDevice.id).then( result =>{
                this.hideLoading(result)

            })
        },
        hideLoading(result){
            this.$q.loading.hide()
            if(result!="Error"){
                this.deviceCharac = result
                this.updateSignal()
            } else {
                this.$q.notify("Couldn't connect. Try Again")
                this.$router.go(-1)
            }
        },
        updateSignal(){
            if(this.signal!="Error"){
                setTimeout(()=>{
                    this.getRSSI()
                },1000)
            } else {
                this.goBack()
            }
        },
        canReadCharac(charac){
            return (charac.Properties.indexOf("Read")>-1)
        },
        readCharac(charac){
            var c = {service: charac.Service, characId: charac.Characteristic}
            psgo.call("web.read.bluetooth.devices", c).then( result =>{
                if(result=="disconnected"){
                    this.goBack()
                } else if(result!="" && result!=null){
                    this.$router.push({name:"info", params:{info:result}})
                } else {
                    this.$q.notify("Nothing to read")
                }
            })     
        },
        goBack(){
            this.$q.notify("Device got disconnected")
            this.$router.go(-1)
        },
        prepareMessage(charac){
            console.log("preparing messgae to " + charac.Service)
            this.$q.dialog({
                title: 'Sending message to ' + this.connectedDevice.title + " - Service " + charac.Service,
                message: 'Enter Message',
                prompt: {
                    model: '',
                    type: 'text'
                },
                ok: {
                    push: true
                },
                cancel: {
                    color: 'negative'
                },
                persistent: true
            }).onOk((m) => {
                psgo.pub("web.write.bluetooth.devices", {service: charac.Service, characId: charac.Characteristic, message: m})
            })
        },
        getRSSI(){
            psgo.call("web.rssi.get.bluetooth.devices",null).then(rssi => {
                this.signal=rssi
                this.updateSignal()
            })
        }
    },
    created(){
        if(this.connectedDevice.id==null){
            psgo.call("web.get.connected.bluetooth.devices", null).then( result =>{
                this.connectedDevice={title:result[0], id:result[1]}
                this.connectToDevice()
            })
        } else {
            this.connectToDevice()            
        }
    },
    beforeDestroy(){
        if(this.$router.currentRoute.name==null){
            psgo.pub("web.disconnect.bluetooth.devices")
        }
        if(this.$q.loading.isActive){
            this.$q.loading.hide()
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
