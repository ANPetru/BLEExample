<template>
    <q-page>
        <div v-if="info!=null">
            <q-list v-if="isArray()" dense  padding >
                <q-item v-for="(key,index) in info" :key="index">
                    <q-item-section>
                        {{index}} - {{key}}
                    </q-item-section>
                </q-item>
            </q-list>
            <div v-else> {{info[0]}} </div>
        </div>
        
    </q-page>
</template>

<script>
export default {
    data(){
        return{
            info: this.$route.params.info,
        }
    },
    methods:{
        isArray(){
            return this.info.length>1
        }
    }
}
</script>

<style>

</style>
