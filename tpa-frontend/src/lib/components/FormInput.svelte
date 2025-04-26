<script lang="ts">
  export let id: string;
  export let type: string = "text";
  export let value: string;
  export let placeholder: string = "";
  export let label: string;
  export let error: string = "";
  export let maxlength: number | undefined = undefined;
  export let pattern: string | undefined = undefined;
  export let showCharCount: boolean = false;

  // Events
  export let onInput: (event: Event) => void = () => {};
  export let onBlur: (event: Event) => void = () => {};
</script>

<div class="form-group">
  <label for={id}>{label}</label>
  <input
    {type}
    {id}
    bind:value
    on:input={onInput}
    on:blur={onBlur}
    class:error
    {placeholder}
    {maxlength}
    {pattern}
  />
  {#if error}
    <p class="error-message">{error}</p>
  {/if}
  {#if showCharCount}
    <p class="char-count">{value.length} / {maxlength || 50}</p>
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

  .char-count {
    text-align: right;
    font-size: 13px;
    color: #6e767d;
    margin-top: 4px;
  }
</style>