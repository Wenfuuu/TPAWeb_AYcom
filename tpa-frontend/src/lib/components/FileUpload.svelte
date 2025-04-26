<script lang="ts">
  export let id: string;
  export let label: string;
  export let fileName: string = "Choose file";
  export let previewUrl: string = '';
  export let error: string = '';
  export let isBanner: boolean = false;
  
  export let onFileChange: (event: Event) => void;
</script>

<div class="form-group">
  <label for={id}>{label}</label>
  <div class="file-upload">
    <input 
      type="file" 
      {id} 
      accept="image/*" 
      on:change={onFileChange}
    />
    {#if previewUrl}
        <div class="image-preview" class:banner-preview={isBanner}>
        <img src={previewUrl} alt="{isBanner ? 'Banner' : 'Profile'} Preview" />
        </div>
    {/if}
    <label for={id} class="file-upload-label">
      {fileName}
    </label>
  </div>
  {#if error}
    <p class="error-message">{error}</p>
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
  
  .file-upload {
    position: relative;
    overflow: hidden;
    display: inline-block;
    width: 100%;
  }
  
  .file-upload input[type="file"] {
    position: absolute;
    left: 0;
    top: 0;
    opacity: 0;
    width: 100%;
    height: 100%;
    cursor: pointer;
  }
  
  .file-upload-label {
    display: block;
    padding: 12px;
    border-radius: 4px;
    border: 1px solid #2f3336;
    background-color: #000;
    color: #fff;
    font-size: 16px;
    text-align: center;
    cursor: pointer;
    overflow: hidden;
    text-overflow: ellipsis;
    white-space: nowrap;
  }
  
  .image-preview {
    margin-top: 12px;
    width: 100px;
    height: 100px;
    border-radius: 50%;
    overflow: hidden;
  }
  
  .banner-preview {
    width: 100%;
    height: 80px;
    border-radius: 4px;
  }
  
  .image-preview img {
    width: 100%;
    height: 100%;
    object-fit: cover;
  }
  
  .error-message {
    color: #f4212e;
    font-size: 13px;
    margin-top: 4px;
  }
</style>