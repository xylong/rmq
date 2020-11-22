<template>
    <div style="text-align: center;margin:100px auto">
        <div style="width:300px;height:200px;border:solid 1px darkgreen;margin: 0 auto;" v-show="loading" >
            正在下单中。。。。
        </div>
        <div><button @click="createOrder">点这里下单</button></div>

    </div>
</template>
<script>
   module.exports ={
       data(){
           return {
                loading: false
           }
       },
       methods:{
            async getResult(no){
                try{
                    const rsp = await  axios.get( 'http://127.0.0.1:8080/result?no=' + no)
                    const { result } = rsp.data
                    if(result > 0 ){
                        alert("下单成功")
                        this.loading = false
                    }else{
                        console.log("继续轮询")
                        setTimeout(()=>this.getResult(no),3000)
                    }
                }catch (e) {
                    console.log(e)
                }
            },
            async createOrder(){
                try{
                    const rsp = await  axios.post( 'http://127.0.0.1:8080/')
                    const { no } = rsp.data
                    this.loading = true
                    console.log(no)
                    this.getResult(no)
                }catch (e) {
                    alert(e)
                }

            }
       }

   }
</script>