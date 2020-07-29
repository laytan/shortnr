<template>
  <div class="d-flex justify-content-center align-items-center vh-100">
    <div class="container">
      <h1 class="d-inline-block border-bottom border-primary">Shortnr</h1>
      <p class="lead w-75">
        Lorem ipsum dolor sit amet, consectetur adipisicing elit.
        Amet, reiciendis.
        Vero assumenda sed non ipsum optio voluptatum
        natus labore sint nostrum tempora enim deleniti aspernatur,
        odio sequi nulla vitae? Placeat reprehenderit deleniti debitis
        natus labore eos tenetur earum mollitia, exercitationem repellat,
        dolores libero quidem.
        Velit tempore delectus odio.
        Reiciendis, quidem.
      </p>
      <div>
        <router-link to="/register" class="btn btn-primary mr-2">Register for free</router-link>
        <router-link to="/login" class="btn btn-secondary">Log in</router-link>
      </div>
    </div>
  </div>
</template>

<script>
export default {
  name: 'Home',
};
</script>

<style lang="scss" scoped>
h1 {
  border-width: 0.25rem !important;
}
</style>
