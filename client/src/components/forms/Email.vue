<template>
  <div>
    <label for="email" class="form-label sr-only">Email</label>
    <div class="input-group">
      <span class="input-group-text bg-primary text-white">
        <Icon icon="envelope" />
      </span>
      <input
        @input="$emit('changed', $event.target.value)"
        autocomplete="email"
        placeholder="Email"
        type="email"
        class="form-control"
        id="email">
    </div>
  </div>
</template>
