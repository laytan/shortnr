<template>
  <div class="grid">
    <div class="bg-secondary background"></div>
    <div class="m-3 bg-primary text-light content d-flex align-items-center">
      <div class="container">
        <h1 class="d-inline-block text-secondary font-weight-bold">Shortnr</h1>
        <p class="lead">
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
          <router-link to="/register" class="btn btn-secondary mr-2">Register for free</router-link>
          <router-link to="/login" class="btn btn-link text-secondary font-weight-bold">
            Log in
          </router-link>
        </div>
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
.grid {
  display: grid;

  .background {
    grid-row: 2 / 10;
    grid-column: 1;
  }

  .content {
    grid-row: 1 / 10;
    grid-column: 1;
  }
}

@media screen and(min-width: 800px) {
  .lead {
    width: 60%;
  }
}
</style>
