<template>
   <select v-on:change="change" multiple v-model="groupsselected">
    <option v-for="(s, value) in groups" :key="value" :value="value">
        {{value}}
    </option>
  </select>
</template>

<script>
import { users as api } from "@/api";

export default {
  name: "groups",
  props: ["group"],
  data: () => {
    return {
      groups:  {},
      groupsselected : {},
    };
  },
  async created() {
    try {
      this.groups = await api.getGroups();
    } catch (e) {
      this.error = e;
    }
  },
  methods: {
    change(event) {
      console.log(event.target.value);
      this.$emit("update:group", this.groupsselected );
    },
  },
};
</script>
