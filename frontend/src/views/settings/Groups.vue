<template>
  <errors v-if="error" :errorCode="error.message" />
  <div class="row" v-else-if="!loading">
    <div class="column" style="width:100%">
      <div class="card">
        <div class="card-title">
          <h2>{{ $t("sidebar.groupManagement") }}</h2>

          <div  class="mobiledisplaynone card-content" style="padding-top:15px;padding-right:30px">
            <input type="checkbox"  v-model="primarygroupvisible" />
            {{ $t("settings.viewuserprimarygroupvisible") }}
          </div>

        <button
          @click="$store.commit('showHover', 'newGroup')"
          class="button">
         {{ $t("buttons.new") }}
        </button>

        </div>

        <div class="card-content full">
          <table>
            <tr>
              <th>Group ID</th>
              <th>GID</th>
              <th class="mobiledisplaynone">Group Members</th>
              <th></th>
            </tr>

            <tr v-for="group in groups" :key="group.id">
                <td v-show=" group.primary==='P' ?  primarygroupvisible : true">{{ group.id }}</td>
                <td v-show=" group.primary==='P' ?  primarygroupvisible : true">{{ group.gid }}</td>
                <td class="mobiledisplaynone" v-show=" group.primary==='P' ?  primarygroupvisible : true">{{ group.members }}</td>
                <td class="small" v-show=" group.primary==='P' ?  primarygroupvisible : true">
                     <i v-if="group.primary==='P'" >&nbsp;</i>
                     <i  v-else class="material-icons"   @click="deleteLink(group.id)">mode_delete</i>
                </td>
            </tr>

          </table>
        </div>
      </div>
    </div>

      <div v-if="$store.state.show === 'deleteGroup'" class="card floating">
      <div class="card-content">
        <p>Are you sure you want to delete this group?</p>
      </div>

      <div class="card-action">
        <button
          class="button button--flat button--grey"
          @click="closeHovers"
          v-focus
          :aria-label="$t('buttons.cancel')"
          :title="$t('buttons.cancel')"
        >
          {{ $t("buttons.cancel") }}
        </button>
        <button class="button button--flat" @click="deleteGroupProcess">
          {{ $t("buttons.delete") }}
        </button>
      </div>
    </div>

  </div>
</template>

<script>
import { mapState, mapMutations } from "vuex";
import { groups as api } from "@/api";
import Errors from "@/views/Errors";
import { BUS } from '@/utils/eventbus';


export default {
  name: "groups",
  components: {
    Errors,
  },
  data: function () {
    return {
      error: null,
      groups: [],
      groupid:"",
      primarygroupvisible:false,
    };
  },
mounted() {
 BUS.$on('bus:groupaddapply', (payload)=>{
        this.groups = [];
        var data = JSON.parse(payload);
        for (let i = 0; i < data.length; i++) {
          this.groups.push(data[i]);
        }
        this.$showSuccess(this.$t("settings.groupCreated"));
      });
    },
  async created() {
    this.setLoading(true);
    try {
      this.groups = await api.getAll();
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
    ...mapMutations(["setLoading","showHover","closeHovers"]),
    deleteLink (idvalue){
      this.groupid = idvalue;
      this.showHover("deleteGroup");
    },
    async deleteGroupProcess(){
       
        try {
             var result =  await api.remove(this.groupid);
        if(result.RESULT_CODE ==="S"){
             var data = await api.getAll();
              this.groups = [];
                for (let i = 0; i < data.length; i++) {
                  this.groups.push(data[i]);
                }
             this.$showSuccess(this.$t("settings.groupDeleted"));
             this.closeHovers();
        }else  if(result.RESULT_CODE ==="F"){
           this.closeHovers();
           this.$showError(result.RESULT_MSG);
        }
      } catch (e) {
        this.$showError(e);
      }
    },
  },
};
</script>

<style scoped>
    @media (max-width:1024px) {
        .mobiledisplaynone {display:none}
    }
</style>