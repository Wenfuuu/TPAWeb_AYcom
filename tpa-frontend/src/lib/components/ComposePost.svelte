<script lang="ts">
  import { type User } from "../model/user";
  import { uploadImage } from "../utils/supabaseStorage";
  import axios from "axios";
  import { fetchThreadCategory, type ThreadCategory } from "../model/thread";
  import { onMount } from "svelte";
  import { fetchCommunity, type Community } from "../model/community";

  export let user: User | null;
  let threadCategories: ThreadCategory[] = [];
  let communities: Community[] = [];
  export let refreshData: () => void = () => {};

  let newPostText = "";
  let uploadedImages: string[] = [];
  let imageFiles: File[] = [];
  let supabaseURLS: string[] = [];
  let currentImageIndex = 0;

  // Variables for category and community selection
  let selectedCategoryId = "";
  let isPostingToCommunity = false;
  let selectedCommunityId = "";
  let showCategoryDropdown = false;
  let showCommunityDropdown = false;

  type ReplySettingType = "everyone" | "following" | "verified";
  let replySettings: ReplySettingType = "everyone";
  let showReplyOptions = false;

  onMount(async () => {
    threadCategories = await fetchThreadCategory();
    communities = await fetchCommunity();
  });

  async function handlePost(): Promise<void> {
    if (newPostText.trim()) {
      if (!selectedCategoryId) {
        alert("Please select a category for your thread");
        return;
      }

      if (isPostingToCommunity && !selectedCommunityId) {
        alert("Please select a community to post to");
        return;
      }

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

      await axios.post("http://localhost:8080/api/threads/post", {
        userId: user?.id,
        categoryId: selectedCategoryId,
        communityId: selectedCommunityId,
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
      selectedCategoryId = "";
      selectedCommunityId = "";
      isPostingToCommunity = false;

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

  function toggleReplyOptions(): void {
    showReplyOptions = !showReplyOptions;
  }

  function setReplyOption(option: ReplySettingType): void {
    replySettings = option;
    showReplyOptions = false;
  }

  // Functions for category and community selection
  function toggleCategoryDropdown(): void {
    showCategoryDropdown = !showCategoryDropdown;
    // Close other dropdowns
    showReplyOptions = false;
    showCommunityDropdown = false;
  }

  function selectCategory(categoryId: string): void {
    selectedCategoryId = categoryId;
    showCategoryDropdown = false;
  }

  function toggleCommunityDropdown(): void {
    showCommunityDropdown = !showCommunityDropdown;
    // Close other dropdowns
    showReplyOptions = false;
    showCategoryDropdown = false;
  }

  function selectCommunity(communityId: string): void {
    selectedCommunityId = communityId;
    showCommunityDropdown = false;
  }

  function togglePostingType(): void {
    isPostingToCommunity = !isPostingToCommunity;
    if (!isPostingToCommunity) {
      selectedCommunityId = "";
    }
  }
</script>

<div class="compose-post">
  <div class="compose-avatar">
    <img src={user?.profile_picture} alt="Your avatar" />
  </div>
  <div class="compose-input">
    <textarea placeholder="What's happening?" bind:value={newPostText}
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

    <!-- Post options section -->
    <div class="post-options">
      <!-- Category selection -->
      <div class="option-dropdown">
        <button
          type="button"
          class="option-button"
          on:click={toggleCategoryDropdown}
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
            <path
              d="M4 20h16a2 2 0 0 0 2-2V8a2 2 0 0 0-2-2h-7.93a2 2 0 0 1-1.66-.9l-.82-1.2A2 2 0 0 0 7.93 3H4a2 2 0 0 0-2 2v13c0 1.1.9 2 2 2Z"
            ></path>
          </svg>
          {selectedCategoryId
            ? threadCategories.find((c) => c.id === selectedCategoryId)?.name ||
              "Select category"
            : "Select category"}
        </button>

        {#if showCategoryDropdown}
          <div class="dropdown-menu">
            {#each threadCategories as category (category.id)}
              <button
                type="button"
                class="dropdown-item"
                on:click={() => selectCategory(category.id)}
              >
                <div class="item-text">
                  <h4>{category.name}</h4>
                </div>
                {#if selectedCategoryId === category.id}
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
                    <polyline points="20 6 9 17 4 12"></polyline>
                  </svg>
                {/if}
              </button>
            {/each}
          </div>
        {/if}
      </div>

      <!-- Post type toggle -->
      <div class="post-type-toggle">
        <button
          type="button"
          class="option-button"
          on:click={togglePostingType}
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
            <path d="M17 21v-2a4 4 0 0 0-4-4H5a4 4 0 0 0-4 4v2"></path>
            <circle cx="9" cy="7" r="4"></circle>
            <path d="M23 21v-2a4 4 0 0 0-3-3.87"></path>
            <path d="M16 3.13a4 4 0 0 1 0 7.75"></path>
          </svg>
          {isPostingToCommunity ? "Posting to community" : "Posting personally"}
        </button>
      </div>

      <!-- Community selection (only shown when posting to community) -->
      {#if isPostingToCommunity}
        <div class="option-dropdown">
          <button
            type="button"
            class="option-button"
            on:click={toggleCommunityDropdown}
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
              <path
                d="M3 5a2 2 0 0 1 2-2h14a2 2 0 0 1 2 2v14a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2V5Z"
              ></path>
              <path d="M3 10h18"></path>
              <path d="M10 3v18"></path>
            </svg>
            {selectedCommunityId
              ? communities.find((c) => c.id === selectedCommunityId)?.name ||
                "Select community"
              : "Select community"}
          </button>

          {#if showCommunityDropdown}
            <div class="dropdown-menu">
              {#each communities as community (community.id)}
                <button
                  type="button"
                  class="dropdown-item"
                  on:click={() => selectCommunity(community.id)}
                >
                  <div class="item-text">
                    <h4>{community.name}</h4>
                    <p>{community.description}</p>
                  </div>
                  {#if selectedCommunityId === community.id}
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
                      <polyline points="20 6 9 17 4 12"></polyline>
                    </svg>
                  {/if}
                </button>
              {/each}
            </div>
          {/if}
        </div>
      {/if}
    </div>

    <div class="reply-settings">
      <button type="button" class="reply-button" on:click={toggleReplyOptions}>
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
          <circle cx="12" cy="12" r="10"></circle>
          <path d="M12 8h.01"></path>
          <path d="M11 12h1v4h1"></path>
        </svg>
        {#if replySettings === "everyone"}
          Everyone can reply
        {:else if replySettings === "following"}
          People you follow can reply
        {:else}
          Only people you mention can reply
        {/if}
      </button>

      {#if showReplyOptions}
        <div class="reply-dropdown">
          <button
            type="button"
            class="reply-option"
            on:click={() => setReplyOption("everyone")}
          >
            <div
              class="option-icon everyone"
              class:selected={replySettings === "everyone"}
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
                <circle cx="12" cy="12" r="10"></circle>
              </svg>
            </div>
            <div class="option-text">
              <h4>Everyone</h4>
              <p>Anyone can reply to this post</p>
            </div>
          </button>

          <button
            type="button"
            class="reply-option"
            on:click={() => setReplyOption("following")}
          >
            <div
              class="option-icon following"
              class:selected={replySettings === "following"}
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
                <path d="M16 21v-2a4 4 0 0 0-4-4H6a4 4 0 0 0-4 4v2"></path>
                <circle cx="9" cy="7" r="4"></circle>
                <path d="M22 21v-2a4 4 0 0 0-3-3.87"></path>
                <path d="M16 3.13a4 4 0 0 1 0 7.75"></path>
              </svg>
            </div>
            <div class="option-text">
              <h4>People you follow</h4>
              <p>Only people you follow can reply</p>
            </div>
          </button>

          <button
            type="button"
            class="reply-option"
            on:click={() => setReplyOption("verified")}
          >
            <div
              class="option-icon mentioned"
              class:selected={replySettings === "verified"}
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
                <path d="M22 12h-4l-3 9L9 3l-3 9H2"></path>
              </svg>
            </div>
            <div class="option-text">
              <h4>Verified only</h4>
              <p>Only verified accounts can reply</p>
            </div>
          </button>
        </div>
      {/if}
    </div>

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
        Post
      </button>
    </div>
  </div>
</div>

<style>
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

  /* Post options */
  .post-options {
    display: flex;
    flex-wrap: wrap;
    gap: 8px;
    margin: 12px 0;
  }

  .option-dropdown {
    position: relative;
  }

  .option-button {
    display: flex;
    align-items: center;
    gap: 8px;
    background: none;
    border: 1px solid var(--border-color);
    border-radius: 16px;
    color: var(--button-primary);
    font-size: 14px;
    padding: 6px 12px;
    cursor: pointer;
  }

  .dropdown-menu {
    position: absolute;
    top: 100%;
    left: 0;
    width: 100%;
    max-height: 200px; /* Set a fixed height for the dropdown */
    overflow-y: auto; /* Enable vertical scrolling */
    background-color: var(--card-background);
    border-radius: 16px;
    box-shadow: 0 0 10px rgba(0, 0, 0, 0.2);
    z-index: 100;
    margin-top: 8px;
    padding: 12px;
    border: 1px solid var(--border-color);
  }

  .dropdown-item {
    width: 100%;
    background-color: var(--card-background);
    border: none;
    display: flex;
    align-items: center;
    justify-content: space-between;
    gap: 12px;
    padding: 12px;
    cursor: pointer;
    border-radius: 8px;
    transition: background-color 0.2s;
    text-align: left;
  }

  .dropdown-item:hover {
    background-color: var(--hover-background);
  }

  .item-text h4 {
    margin: 0;
    font-size: 15px;
    font-weight: 700;
    color: var(--text-primary);
  }

  .item-text p {
    margin: 4px 0 0 0;
    font-size: 13px;
    color: var(--text-tertiary);
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

  /* Reply settings styles */
  .reply-settings {
    position: relative;
    margin: 12px 0;
  }

  .reply-button {
    display: flex;
    align-items: center;
    gap: 8px;
    background: none;
    border: none;
    color: var(--button-primary);
    font-size: 14px;
    cursor: pointer;
    padding: 4px 0;
  }

  .reply-dropdown {
    position: absolute;
    top: 100%;
    left: 0;
    width: 300px;
    background-color: var(--card-background);
    border-radius: 16px;
    box-shadow: 0 0 10px rgba(0, 0, 0, 0.2);
    z-index: 100;
    margin-top: 8px;
    padding: 12px;
    border: 1px solid var(--border-color);
  }

  .reply-option {
    width: 100%;
    background-color: var(--background);
    border: none;
    display: flex;
    align-items: flex-start;
    gap: 12px;
    padding: 12px;
    cursor: pointer;
    border-radius: 8px;
    transition: background-color 0.2s;
    text-align: left;
  }

  .reply-option:hover {
    background-color: var(--hover-background);
  }

  .option-icon {
    width: 40px;
    height: 40px;
    border-radius: 50%;
    background-color: var(--button-primary);
    display: flex;
    align-items: center;
    justify-content: center;
    color: white;
  }

  .option-icon.selected {
    background-color: var(--button-primary);
  }

  .option-text h4 {
    margin: 0;
    font-size: 15px;
    font-weight: 700;
    color: var(--text-primary);
  }

  .option-text p {
    margin: 4px 0 0 0;
    font-size: 13px;
    color: var(--text-tertiary);
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
