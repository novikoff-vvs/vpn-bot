<template>
  <div class="container">
    <div class="row">
      <div class="col-1">
        <div class="card" v-on:click="goToProfile">
          <h2><BsPersonArmsUp/> Профиль</h2>
          <p>Информация о вашем аккаунте</p>
        </div>
      </div>
    </div>
    <div class="row">
      <div class="col-1">
        <div class="card border border-light" v-on:click="goToPayments">
          <h2><FaDollarSign/> Оплата </h2>
          <p>История платежей и система оплаты</p>
        </div>
      </div>
    </div>
  </div>

</template>
<script setup>
import { useRouter } from 'vue-router'
import { BsPersonArmsUp,FaDollarSign } from '@kalimahapps/vue-icons';

const router = useRouter()

const goToPayments = () => {
  router.push('/payments')
}
const goToProfile = () => {
  router.push('/profile')
}
</script>