<script setup lang="ts">
// This starter template is using Vue 3 <script setup> SFCs
// Check out https://vuejs.org/api/sfc-script-setup.html#script-setup
import { AddressChecker, sm } from './address'
import { ref } from 'vue'
import { data } from './store'
const check = ref(new AddressChecker('http://localhost:8073/address', data))
const compare = ref<string[][]>([])

const snap = ref(data)


function modify() {
  data.set("PRV-CITY", "MARPLE")
  data.set("PRT-CITY", "BROOMALL")
  console.log(sm(data), sm(check.value.old))
}
async function all() {
  compare.value = await check.value.suggest(check.value.addresses)
}
async function changed() {
  compare.value = await check.value.suggest(check.value.changes)
}
function replace() {
  let changes = check.value.accept()
  console.log("changes", sm(changes))
  snap.value = data
}


</script>

<template>
  <div>
    <button @click="modify">modify</button>
    <button @click="all">all</button>
    <button @click="changed">changed</button>
    <button @replace="replace">replace</button>

    <h2>Suggestions</h2>
    <table>
      <tr v-for="(item,index) in compare">
        <td>{{ item[0]}}</td>
        <td>{{ item[1]}}</td>
      </tr>
    </table>
    <hr />
    <button>cancel</button><button>accept</button>
    <h2> data </h2>
    <pre>{{ sm(data) }}</pre>


  </div>
</template>

