<template>
    <q-page padding>
        <div v-if="deviceCharac==''">
        Connecting To: {{connectedDevice.title}}<br/>
        </div>
        <div v-else>
            Connected To : {{connectedDevice.title}}<br/>
            <q-list bordered separator>
                <q-item v-for="(charac,index) in deviceCharac" :key="index" class="row">
                    <q-item-section class="col">Service: {{ charac.Service }} - {{charac.Characteristic}}
                    </q-item-section>
                    <q-item-section side>
                        <q-icon class="col" v-if="canReadCharac(charac)" name="visibility" @click.native="readCharac(charac)"></q-icon>
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
        }
    }, 
    methods:{
        connectToDevice(){
            psgo.call("web.connect.bluetooth.devices", this.connectedDevice.id).then( result =>{
                if(result!="Error"){
                    this.deviceCharac = result
                } else {
                    this.$q.notify("Couldn't connect. Try Again")
                    this.$router.go(-1)
                }
            })
        },
        canReadCharac(charac){
            return (charac.Properties.indexOf("Read")>-1)
        },
        readCharac(charac){
            if(this.canReadCharac(charac)){
                var c = {service: charac.Service, characId: charac.Characteristic}
                psgo.call("web.read.bluetooth.devices", c).then( result =>{
                    if(result!="" && result!=null){
                        this.$router.push({name:"info", params:{info:result}})
                    } else {
                        this.$q.notify("Nothing to read")
                    }
                })
            }
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
    }
}
</script>

<style>

</style>
