<template>
  <div>
    <nav class="navbar navbar-expand-md navbar-dark bg-secondary">
      <div class="container-fluid">
        <router-link class="navbar-brand text-primary" to="/">Shortnr</router-link>
        <button
          class="navbar-toggler"
          data-toggle="collapse"
          data-target="#primaryNav"
          aria-controls="primaryNav"
          aria-expanded="false"
          aria-label="Toggle navigation"
        >
          <span class="navbar-toggler-icon"></span>
        </button>
        <div class="collapse navbar-collapse" id="primaryNav">
          <ul class="navbar-nav ml-auto mb-2 mb-md-0">
            <li class="nav-item">
              <router-link class="nav-link" to="/">Login</router-link>
            </li>
            <li class="nav-item">
              <router-link class="nav-link" to="/register">Register</router-link>
            </li>
          </ul>
        </div>
      </div>
    </nav>
    <router-view/>
  </div>
</template>
