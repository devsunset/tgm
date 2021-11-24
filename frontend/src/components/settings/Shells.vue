<template>
   <select  ref="xshell" v-on:change="change">
    <option v-for="(s, value) in shells" :key="value" :value="value" >
        {{value}}
    </option>
  </select>
</template>

<script>
import { users as api } from "@/api";

export default {
  name: "shells",
  props: ["shell"],
  data: () => {
    return {
      shells:  {},
    };
  },
  async created() {
    try {
      this.shells = await api.getShells();
    } catch (e) {
      this.error = e;
    }
  },
  updated: function () {
  this.$nextTick(function () {
      this.$refs.xshell.value = this.shell;
  })
},
  methods: {
    change(event) {
      this.$emit("update:shell", event.target.value);
    },
  },
};
</script>
