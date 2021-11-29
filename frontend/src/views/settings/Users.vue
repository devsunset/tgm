<template>
  <errors v-if="error" :errorCode="error.message" />
  <div class="row" v-else-if="!loading">
    <div class="column" style="width:100%">
      <div class="card">
        <div class="card-title">
          <h2>{{ $t("settings.userManagement") }}</h2>
          <label for="username" class="mobiledisplaynone" style="padding-top:15px;padding-right:10px">{{ $t("settings.username") }}</label>
            <input
              class="mobiledisplaynone input"
              type="text"
              v-model="username"
              @keyup.enter="userSearch"
              id="username"
            />
            <div class="mobiledisplaynone" style="padding-right:20px">
            <button class="button" @click="userSearch">
              {{ $t("buttons.search") }}
            </button>
            </div>
          <router-link to="/settings/users/new"
            ><button class="button" >
              {{ $t("buttons.new") }}
            </button></router-link
          >
        </div>

        <div class="card-content full">
          <table>
            <tr>
              <th>{{ $t("settings.admin") }}</th>
              <th>{{ $t("settings.username") }}</th>
              <th>{{ $t("settings.shell") }}</th>
              <th class="mobiledisplaynone">{{ $t("settings.group") }}</th>
              <th class="mobiledisplaynone">{{ $t("settings.scope") }}</th>
              <th></th>
              <th></th>
            </tr>

            <tr v-for="user in users" :key="user.id">
              <td>
                <i v-if="user.perm.admin" class="material-icons">person</i
                ><i v-else>&nbsp;</i>
              </td>
              <td>{{ user.username }}</td>
              <td>{{ user.shell }}</td>
              <td class="mobiledisplaynone">{{ user.group }}</td>
              <td class="mobiledisplaynone">{{ user.scope }}</td>
              <td class="small">
                    <i v-if="user.lock === 'LK'" class="material-icons">locked</i
                    ><i v-else-if="user.lock === 'NP'" class="material-icons">lock_outline</i>
                    <i v-else>&nbsp;</i>
                </td>
              <td class="small">
                <router-link :to="'/settings/users/' + user.id"
                  ><i class="material-icons">mode_edit</i></router-link
                >
              </td>
            </tr>
          </table>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import { mapState, mapMutations } from "vuex";
import { users as api } from "@/api";
import Errors from "@/views/Errors";

export default {
  name: "users",
  components: {
    Errors,
  },
  data: function () {
    return {
      error: null,
      users: [],
      username: "",
    };
  },
  async created() {
    this.setLoading(true);

    try {
      this.users = await api.getAll("");
    } catch (e) {
      this.error = e;
    } finally {
      this.setLoading(false);
    }
  },
  computed: {
    ...mapState(["loading"]),
  },
  methods: {
    ...mapMutations(["setLoading"]),
    userSearch: async function () {
      //this.setLoading(true);
      try {
        this.groups = [];
        this.users = await api.getAll(this.username);
      } catch (e) {
        this.error = e;
      } 
      //finally {
      //  this.setLoading(false);
      //}
    }
  },
};
</script>

<style scoped>
    @media (max-width:1024px) {
        .mobiledisplaynone {display:none}
    }
</style>