<script lang="ts">
  import { fetchThreadDetailWithUser, type Thread } from "../model/thread";
  import { type User, type UserThread } from "../model/user";
  import { uploadImage } from "../utils/supabaseStorage";
  import axios from "axios";
  import { onMount } from "svelte";

  export let thread: Thread & { userData: UserThread };
  export let user: User | null;
  export let refreshData: () => void = () => {};

  let newPostText = "";
  let uploadedImages: string[] = [];
  let imageFiles: File[] = [];
  let supabaseURLS: string[] = [];
  let currentImageIndex = 0;

  type ReplySettingType = "everyone" | "following" | "verified";
  const replySettings: ReplySettingType = "everyone";

  onMount(async () => {
    const currentUrl = window.location.href;
    const urlParts = currentUrl.split("/");
    const threadId = urlParts[urlParts.length - 1];

    if (user) {
      const data = await fetchThreadDetailWithUser(threadId, user.id);
      if (data) thread = data;
    }
  });

  async function handlePost(): Promise<void> {
    if (newPostText.trim()) {
      console.log("Posting:", newPostText, uploadedImages);

      // upload images to supabase
      if (imageFiles.length > 0) {
        await uploadImagesToSupabase();
      }

      // call api to create thread
      let isAdvertisement;
      if (user?.role == "user") isAdvertisement = false;
      else isAdvertisement = true;
      // console.log(user);

      await axios.post("http://localhost:8080/api/threads/reply", {
        userId: user?.id,
        replyToId: thread.id,
        content: newPostText,
        replyPermission: replySettings,
        isAdvertisement,
        scheduledAt: "",
        imageUrls: supabaseURLS,
      });

      newPostText = "";
      uploadedImages = [];
      imageFiles = [];
      supabaseURLS = [];

      refreshData();
    }
  }

  // upload and get urls
  async function uploadImagesToSupabase(): Promise<boolean> {
    if (!imageFiles || imageFiles.length === 0) return false;

    // upload all image here
    for (const file of imageFiles) {
      const supaURL = await uploadImage(file);
      if (supaURL) {
        supabaseURLS.push(supaURL);
      }
    }

    return true;
  }

  // for image preview
  function handleImageUpload(event: Event): void {
    const target = event.target as HTMLInputElement;
    const files = target.files;

    if (!files || files.length === 0) return;

    for (let i = 0; i < files.length; i++) {
      const file = files[i];
      imageFiles = [...imageFiles, file];
      const reader = new FileReader();

      reader.onload = (e: ProgressEvent<FileReader>) => {
        const result = e.target?.result;
        if (typeof result === "string") {
          uploadedImages = [...uploadedImages, result];
        }
      };

      reader.readAsDataURL(file);
    }
  }

  function removeImage(index: number): void {
    uploadedImages = uploadedImages.filter((_, i) => i !== index);
    imageFiles = imageFiles.filter((_, i) => i !== index);
    if (currentImageIndex >= uploadedImages.length) {
      currentImageIndex = Math.max(0, uploadedImages.length - 1);
    }
  }
</script>

{#if user && thread}
    <div class="reply-to">
      Replying to
      <span class="reply-username"
        >@{thread.userData?.username || "unknown"}</span
      >
    </div>
  <div class="compose-post">
    <div class="compose-avatar">
      <img src={user?.profile_picture} alt="Your avatar" />
    </div>
    <div class="compose-input">
      <textarea placeholder="Post your reply" bind:value={newPostText}
      ></textarea>

      {#if uploadedImages.length > 0}
        <div class="image-carousel">
          <div class="carousel-container">
            <div
              class="carousel-images"
              class:has-multiple={uploadedImages.length > 1}
            >
              {#each uploadedImages as image, i (i)}
                <div
                  class="carousel-image"
                  class:active={i === currentImageIndex}
                >
                  <img
                    src={image || "/placeholder.svg"}
                    alt="Uploaded image {i + 1}"
                  />
                  <button type="button" class="edit-button">Edit</button>
                  <button
                    type="button"
                    class="remove-button"
                    aria-label="Remove Button"
                    on:click={() => removeImage(i)}
                  >
                    <svg
                      xmlns="http://www.w3.org/2000/svg"
                      width="16"
                      height="16"
                      viewBox="0 0 24 24"
                      fill="none"
                      stroke="currentColor"
                      stroke-width="2"
                      stroke-linecap="round"
                      stroke-linejoin="round"
                    >
                      <line x1="18" y1="6" x2="6" y2="18"></line>
                      <line x1="6" y1="6" x2="18" y2="18"></line>
                    </svg>
                  </button>
                </div>
              {/each}
            </div>

            {#if uploadedImages.length > 1}
              <div class="carousel-nav">
                <button
                  type="button"
                  class="nav-button prev"
                  aria-label="Button Prev"
                  disabled={currentImageIndex === 0}
                  on:click={() =>
                    (currentImageIndex = Math.max(0, currentImageIndex - 1))}
                >
                  <svg
                    xmlns="http://www.w3.org/2000/svg"
                    width="20"
                    height="20"
                    viewBox="0 0 24 24"
                    fill="none"
                    stroke="currentColor"
                    stroke-width="2"
                    stroke-linecap="round"
                    stroke-linejoin="round"
                  >
                    <path d="m15 18-6-6 6-6"></path>
                  </svg>
                </button>
                <button
                  type="button"
                  class="nav-button next"
                  aria-label="Button Next"
                  disabled={currentImageIndex === uploadedImages.length - 1}
                  on:click={() =>
                    (currentImageIndex = Math.min(
                      uploadedImages.length - 1,
                      currentImageIndex + 1
                    ))}
                >
                  <svg
                    xmlns="http://www.w3.org/2000/svg"
                    width="20"
                    height="20"
                    viewBox="0 0 24 24"
                    fill="none"
                    stroke="currentColor"
                    stroke-width="2"
                    stroke-linecap="round"
                    stroke-linejoin="round"
                  >
                    <path d="m9 18 6-6-6-6"></path>
                  </svg>
                </button>
              </div>

              <div class="carousel-indicators">
                {#each uploadedImages as _, i (i)}
                  <span class="indicator" class:active={i === currentImageIndex}
                  ></span>
                {/each}
              </div>
            {/if}
          </div>
        </div>
      {/if}

      <div class="compose-actions">
        <div class="compose-tools">
          <!-- Image upload button -->
          <label class="tool-button" for="image-upload">
            <svg
              xmlns="http://www.w3.org/2000/svg"
              width="20"
              height="20"
              viewBox="0 0 24 24"
              fill="none"
              stroke="currentColor"
              stroke-width="2"
              stroke-linecap="round"
              stroke-linejoin="round"
            >
              <rect width="18" height="18" x="3" y="3" rx="2" ry="2"></rect>
              <circle cx="9" cy="9" r="2"></circle>
              <path d="m21 15-3.086-3.086a2 2 0 0 0-2.828 0L6 21"></path>
            </svg>
            <input
              type="file"
              id="image-upload"
              accept="image/*"
              multiple
              on:change={handleImageUpload}
              class="hidden-input"
            />
          </label>

          <!-- Poll button -->
          <button type="button" class="tool-button" aria-label="Create poll">
            <svg
              xmlns="http://www.w3.org/2000/svg"
              width="20"
              height="20"
              viewBox="0 0 24 24"
              fill="none"
              stroke="currentColor"
              stroke-width="2"
              stroke-linecap="round"
              stroke-linejoin="round"
            >
              <path d="M3 3v18h18"></path>
              <path d="M18 17V9"></path>
              <path d="M13 17V5"></path>
              <path d="M8 17v-3"></path>
            </svg>
          </button>

          <!-- Schedule button -->
          <button type="button" class="tool-button" aria-label="Schedule post">
            <svg
              xmlns="http://www.w3.org/2000/svg"
              width="20"
              height="20"
              viewBox="0 0 24 24"
              fill="none"
              stroke="currentColor"
              stroke-width="2"
              stroke-linecap="round"
              stroke-linejoin="round"
            >
              <circle cx="12" cy="12" r="10"></circle>
              <polyline points="12 6 12 12 16 14"></polyline>
            </svg>
          </button>
        </div>

        <button
          type="button"
          class="post-button"
          disabled={!newPostText.trim() && uploadedImages.length === 0}
          on:click={handlePost}
        >
          Reply
        </button>
      </div>
    </div>
  </div>
{/if}

<style>
   .reply-to {
    display: flex;
    align-items: center;
    gap: 6px;
    font-size: 14px;
    color: var(--text-tertiary);
    margin: 0 0 8px 60px;
    padding: 4px 12px;
  }

  .reply-username {
    color: var(--button-primary);
  }

  .compose-post {
    display: flex;
    padding: 16px;
    border-bottom: 1px solid var(--border-color);
  }

  .compose-avatar {
    margin-right: 12px;
  }

  .compose-avatar img {
    width: 48px;
    height: 48px;
    border-radius: 50%;
  }

  .compose-input {
    flex: 1;
    display: flex;
    flex-direction: column;
  }

  textarea {
    background: none;
    border: none;
    color: var(--text-primary);
    font-size: 20px;
    font-family: inherit;
    height: 52px;
    padding: 12px 0;
    resize: none;
    outline: none;
  }

  textarea::placeholder {
    color: var(--text-tertiary);
  }

  /* Hide file input but keep it accessible */
  .hidden-input {
    display: none;
  }

  /* Image carousel styles */
  .image-carousel {
    margin: 12px 0;
    border-radius: 16px;
    overflow: hidden;
    border: 1px solid var(--border-color);
  }

  .carousel-container {
    position: relative;
    width: 100%;
    height: 280px;
    background-color: #000;
    overflow: hidden;
  }

  .carousel-images {
    display: flex;
    width: 100%;
    height: 100%;
  }

  .carousel-images.has-multiple {
    flex-wrap: nowrap;
  }

  .carousel-image {
    position: relative;
    min-width: 100%;
    height: 100%;
    display: none;
  }

  .carousel-image.active {
    display: block;
  }

  .carousel-image img {
    width: 100%;
    height: 100%;
    object-fit: contain;
  }

  .edit-button {
    position: absolute;
    top: 8px;
    right: 8px;
    background-color: rgba(0, 0, 0, 0.6);
    color: white;
    border: none;
    border-radius: 9999px;
    padding: 6px 12px;
    font-size: 14px;
    font-weight: 600;
    cursor: pointer;
  }

  .remove-button {
    position: absolute;
    top: 8px;
    left: 8px;
    background-color: rgba(0, 0, 0, 0.6);
    color: white;
    border: none;
    border-radius: 50%;
    width: 30px;
    height: 30px;
    display: flex;
    align-items: center;
    justify-content: center;
    cursor: pointer;
  }

  .carousel-nav {
    position: absolute;
    top: 0;
    left: 0;
    width: 100%;
    height: 100%;
    pointer-events: none;
  }

  .nav-button {
    position: absolute;
    top: 50%;
    transform: translateY(-50%);
    background-color: rgba(0, 0, 0, 0.6);
    color: white;
    border: none;
    border-radius: 50%;
    width: 40px;
    height: 40px;
    display: flex;
    align-items: center;
    justify-content: center;
    cursor: pointer;
    pointer-events: auto;
  }

  .nav-button:disabled {
    opacity: 0.5;
    cursor: not-allowed;
  }

  .nav-button.prev {
    left: 8px;
  }

  .nav-button.next {
    right: 8px;
  }

  .carousel-indicators {
    position: absolute;
    bottom: 12px;
    left: 50%;
    transform: translateX(-50%);
    display: flex;
    gap: 4px;
  }

  .indicator {
    width: 8px;
    height: 8px;
    border-radius: 50%;
    background-color: rgba(255, 255, 255, 0.5);
  }

  .indicator.active {
    background-color: white;
  }

  .compose-actions {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-top: 12px;
    border-top: 1px solid var(--border-color);
    padding-top: 12px;
  }

  .compose-tools {
    display: flex;
    gap: 4px;
  }

  .tool-button {
    background: none;
    border: none;
    color: var(--button-primary);
    width: 34px;
    height: 34px;
    border-radius: 50%;
    display: flex;
    align-items: center;
    justify-content: center;
    cursor: pointer;
    transition: background-color 0.2s;
  }

  .tool-button:hover {
    background-color: var(--hover-background);
  }

  .post-button {
    background-color: var(--button-primary);
    color: var(--button-primary-text);
    border: none;
    border-radius: 9999px;
    padding: 8px 16px;
    font-size: 15px;
    font-weight: 700;
    cursor: pointer;
    transition: background-color 0.2s;
  }

  .post-button:hover:not(:disabled) {
    background-color: var(--button-primary-hover);
  }

  .post-button:disabled {
    background-color: var(--button-primary-hover);
    opacity: 0.5;
    cursor: not-allowed;
  }
</style>
