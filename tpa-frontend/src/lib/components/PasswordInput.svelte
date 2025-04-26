<script lang="ts">
  import type { PasswordValidations } from "../utils/registerValidation";

  export let id: string;
  export let value: string;
  export let placeholder: string;
  export let label: string;
  export let error: string = "";
  export let showRequirements: boolean = false;
  export let validations: PasswordValidations | null = null;

  export let onInput: (event: Event) => void = () => {};
  export let onBlur: (event: Event) => void = () => {};
</script>

<div class="form-group">
  <label for={id}>{label}</label>
  <input
    type="password"
    {id}
    bind:value
    on:input={onInput}
    on:blur={onBlur}
    class:error
    {placeholder}
  />
  {#if error}
    <p class="error-message">{error}</p>
  {/if}

  {#if showRequirements && validations}
    <div class="password-requirements">
      <p class:valid={validations.length}>
        ✓ At least 8 characters
      </p>
      <p class:valid={validations.uppercase}>
        ✓ At least 1 uppercase letter
      </p>
      <p class:valid={validations.lowercase}>
        ✓ At least 1 lowercase letter
      </p>
      <p class:valid={validations.number}>
        ✓ At least 1 number
      </p>
    </div>
  {/if}
</div>

<style>
  .form-group {
    margin-bottom: 20px;
  }

  label {
    display: block;
    font-size: 14px;
    font-weight: 500;
    margin-bottom: 8px;
    color: #e7e9ea;
  }

  input {
    width: 100%;
    padding: 12px;
    border-radius: 4px;
    border: 1px solid #2f3336;
    background-color: #000;
    color: #fff;
    font-size: 16px;
    box-sizing: border-box;
  }

  input.error {
    border-color: #f4212e;
  }

  .error-message {
    color: #f4212e;
    font-size: 13px;
    margin-top: 4px;
  }

  .password-requirements {
    margin-top: 8px;
    font-size: 13px;
    color: #6e767d;
  }

  .password-requirements p {
    margin: 4px 0;
  }

  .password-requirements p.valid {
    color: #00ba7c;
  }
</style>