<template>
  <errors v-if="error" :errorCode="error.message" />
  <div class="row" v-else-if="!loading">
    <div class="column" style="width:100%">
      <div class="card">
        <div class="card-title">
          <h2>{{ $t("settings.userManagement") }}</h2>

          <label for="username" style="padding-top:15px;padding-right:10px">{{ $t("settings.username") }}</label>
            <input
              :class="userNameClass"
              type="text"
              v-model="username"
              @keyup.enter="userSearch"
              id="username"
            />
            <div style="padding-right:20px">
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
               <th>{{ $t("settings.scope") }}</th>
              <th>Shell</th>
              <th>Group</th>
              <th></th>
              <th></th>
            </tr>

            <tr v-for="user in users" :key="user.id">
              <td>
                <i v-if="user.perm.admin" class="material-icons">person</i
                ><i v-else>&nbsp;</i>
              </td>
              <td>{{ user.username }}</td>
              <td>{{ user.scope }}</td>
              <td>{{ user.shell }}</td>
              <td>{{ user.group }}</td>
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
    },
  },
};
</script>
