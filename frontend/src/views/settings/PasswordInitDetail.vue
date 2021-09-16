<template>
  <div>
    <p>
      <label for="password">{{ $t("settings.password") }}</label>
      <input
        class="input input--block"
        type="password"
        :placeholder="$t('settings.newPassword')"
        v-model="password"
        id="password"
      />
      <input
          class="input input--block"
          type="password"
          :placeholder="$t('settings.newPasswordConfirm')"
          v-model="passwordConf"
          name="password"
      />
    </p>
  </div>
</template>

<script>
export default {
  name: "user",
  data: function () {
    return {   
      password : "",   
      passwordConf: "",
    };
  },
  props: ["user", "isNew"],
  computed: {
    passwordClass() {
      const baseClass = "input input--block";
      if (this.password === ""  && this.passwordConf === "") {
        return baseClass;
      }
             
      var reg = /^(?=.*?[A-Z])(?=.*?[a-z])(?=.*?[0-9])(?=.*?[#?!@$%^&*-]).{8,}$/;
      var hangulcheck = /[ㄱ-ㅎ|ㅏ-ㅣ|가-힣]/;
      
      if(false === reg.test(this.password)) {
        return `${baseClass} input--red`;
      }else if(/(\w)\1\1\1/.test(this.password)){
       return `${baseClass} input--red`;
      }else if(this.password.search(/\s/) != -1){
        return `${baseClass} input--red`;
      }else if(hangulcheck.test(this.password)){
       return `${baseClass} input--red`;
      }else {
        if (this.password !== this.passwordConf || this.password === "") {
          return `${baseClass} input--red`;
        }
      }

      if (this.password === this.passwordConf) {
        return `${baseClass} input--green`;
      }

      return `${baseClass} input--red`;
    },
  },
};
</script>
