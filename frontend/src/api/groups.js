import { fetchURL, fetchJSON } from "./utils";

export async function getAll() {
  return fetchJSON(`/api/groups`, {});
}

export async function create(groupid) {
  const res = await fetchURL(`/api/groups`, {
    method: "POST",
    body: JSON.stringify({
      what: "group",
      which: [],
      data: groupid,
    }),
  });

  if (res.status !== 200) {
    throw new Error(res.status);
  }
}

export async function remove(groupid) {
  const res = await fetchURL(`/api/groups`, {
    method: "DELETE",
    body: JSON.stringify({
      what: "group",
      which: [],
      data: groupid,
    }),
  });

  if (res.status !== 200) {
    throw new Error(res.status);
  }
}
