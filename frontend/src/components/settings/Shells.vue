<template>
   <select v-on:change="change">
    <option v-for="(s, value) in shells" :key="value" :value="value">
        {{value}}
    </option>
  </select>
</template>

<script>
import { users as api } from "@/api";

export default {
  name: "shells",
  props: ["shell"],
  data() {
    let dataObj = {
       shells:  {},
    };

    return dataObj;
  },
async created() {
    try {
      this.shells = await api.getShells();
      // Object.keys(this.shells).forEach(function (key) {
      //   this.$emit("update:shell", key); 
      //   return
      // });

    } catch (e) {
      this.error = e;
    }
  },
  methods: {
    change(event) {
      this.$emit("update:shell", event.target.value);
    },
  },
};
</script>
