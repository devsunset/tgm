<template>
  <select :value="group" multiple="multiple">
    <option v-for="(group, value) in groups" :key="value" :value="value">
      {{ $t("groups." + group) }}
    </option>
  </select>
</template>

<script>
export default {
  name: "groups",
  props: ["group"],
  data() {
    let dataObj = {
       groups: {
        en: "en",
        ko: "ko",
      },
    };

Object.defineProperty(dataObj, "groups", {
      configurable: false,
      writable: false,
    });

    return dataObj;
  },
 
};
</script>
