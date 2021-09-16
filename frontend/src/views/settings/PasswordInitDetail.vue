<template>
  <div>
    <p>
      <label for="password">{{ $t("settings.username") }}</label>
      <input
        class="input input--block input--gray"
        type="text"
        v-model="requestId"
        name=""
        disabled="true"
      />
    </p>
    <p>
      <label for="password">이름</label>
      <input
        class="input input--block input--gray"
        type="text"
        v-model="requestName"
        name=""
        disabled="true"
      />
    </p>
    <p>
      <label for="password">전화번호</label>
      <input
        class="input input--block input--gray"
        type="text"
        v-model="requestPhoneNumber"
        name=""
        disabled="true"
      />
    </p>
    <p>
      <label for="password">Request Date</label>
      <input
        class="input input--block input--gray"
        type="text"
        v-model="requestDate"
        name=""
        disabled="true"
      />
    </p>
    <p>
      <label for="password">Apply Date</label>
      <input
        class="input input--block input--gray"
        type="text"
        v-model="applyDate"
        name=""
        disabled="true"
      />
    </p>
    <p>
      <label for="password">Status</label>
      <input
        class="input input--block input--gray"
        type="text"
        v-model="requestStatus"
        name=""
        disabled="true"
      />
    </p>
    <p>
      <label for="password">{{ $t("settings.password") }}</label>
      <input
        :class="passwordClass"
        type="password"
        :placeholder="$t('settings.newPassword')"
        v-model="password"
        name="password"
      />
      <input
          :class="passwordSubClass"
          type="password"
          :placeholder="$t('settings.newPasswordConfirm')"
          v-model="passwordConf"
          name="passwordConf"
      />
    </p>

    <div class="card-action" style="float:right">
      <button type="button" aria-label="삭제" title="삭제" class="button button--flat button--red"> 삭제 </button>
      <input type="submit" class="button button--flat" value="저장">
    </div>
    
  </div>
  
</template>

<script>
export default {
  name: "user",
  data: function () {
    return {   
      requestId :"tgm1",
      requestName :"홍길동",
      requestPhoneNumber :"010-9999-7777",
      requestDate :"2021-09-19 12:00:00",
      applyDate :"",
      requestStatus :"Request",
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
       return `${baseClass} input--green`;
      }

    },
    passwordSubClass() {
      const baseClass = "input input--block";
      if (this.password === ""  && this.passwordConf === "") {
        return baseClass;
      }
             
      if (this.password === this.passwordConf) {
        return `${baseClass} input--green`;
      }else{
        return `${baseClass} input--red`;
      }

    },
  },
};
</script>
