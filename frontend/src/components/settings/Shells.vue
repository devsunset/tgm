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

// Object.defineProperty(dataObj, "shells", {
//       configurable: false,
//       writable: false,
//     });

    return dataObj;
  },
async created() {
    try {
      this.shells = await api.getShells();
    } catch (e) {
      this.error = e;
    }
  },
  methods: {
    change(event) {
      console.log(event.target.value);
      //this.$emit("update:locale", event.target.value);
    },
  },
};
</script>
