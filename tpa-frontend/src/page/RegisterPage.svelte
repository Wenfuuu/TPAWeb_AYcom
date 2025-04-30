<script module lang="ts">
  declare global {
    interface Window {
      recaptchaCallback: (token: string) => void;
    }
  }
</script>

<script lang="ts">
  import { onMount } from "svelte";
  import { navigate } from "svelte-routing";
  import FileUpload from "../lib/components/FileUpload.svelte";
  import FormInput from "../lib/components/FormInput.svelte";
  import PasswordInput from "../lib/components/PasswordInput.svelte";
  import Button from "../lib/components/Button.svelte";
  import {
    validateName,
    validateUsername,
    validateEmail,
    validatePassword,
    validateConfirmPassword,
    validateGender,
    validateDateOfBirth,
    validateSecurityQuestion,
    validateProfilePicture,
    validateBanner,
    validateVerificationCode,
    type FormData,
    type ErrorMessages,
    type PasswordValidations,
    type MonthOption,
    type DayOption,
    type YearOption,
  } from "../lib/utils/registerValidation";
  import {
    uploadProfilePicture,
    uploadBannerImage,
  } from "../lib/utils/supabaseStorage";
  import { checkGuest } from "../lib/utils/auth";

  // Define the API data structure
  type ApiFormData = {
    name: string;
    username: string;
    email: string;
    password: string;
    gender: string;
    profilePicture: string; // URL to Supabase image
    banner: string; // URL to Supabase image
    dateOfBirth: string;
    isSubscribed: boolean;
    securityQuestion: string;
    securityAnswer: string;
  };

  const securityQuestions = [
    { id: 1, text: "What was the name of your first pet?" },
    { id: 2, text: "What city were you born in?" },
    { id: 3, text: "What is your favorite video game?" },
    { id: 4, text: "What was the name of your first school?" },
    { id: 5, text: "What was your childhood nickname?" },
  ];

  const state = $state({
    currentStep: 1,
    isLoading: false,
    uploadingImages: false,
    verificationCodeSent: false,
    verificationCodeExpiry: null as Date | null,
    resendDisabled: false,
    resendCountdown: 0,
    recaptchaVerified: false,
    formData: {
      name: "",
      username: "",
      email: "",
      password: "",
      confirmPassword: "",
      gender: "",
      profilePicture: null as File | null,
      profilePicturePreview: "",
      banner: null as File | null,
      bannerPreview: "",
      birthDay: "",
      birthMonth: "",
      birthYear: "",
      securityQuestion: "",
      securityAnswer: "",
      subscribeToNewsletter: false,
      verificationCode: "",
    } as FormData,
    // Store the Supabase URLs separately
    profilePictureUrl: "",
    bannerUrl: "",
    errors: {
      name: "",
      username: "",
      email: "",
      password: "",
      confirmPassword: "",
      gender: "",
      profilePicture: "",
      banner: "",
      dateOfBirth: "",
      securityQuestion: "",
      securityAnswer: "",
      verificationCode: "",
      general: "",
    } as ErrorMessages,
    passwordValidations: {
      length: false,
      uppercase: false,
      lowercase: false,
      number: false,
    } as PasswordValidations,
    months: [
      { value: "01", text: "January" },
      { value: "02", text: "February" },
      { value: "03", text: "March" },
      { value: "04", text: "April" },
      { value: "05", text: "May" },
      { value: "06", text: "June" },
      { value: "07", text: "July" },
      { value: "08", text: "August" },
      { value: "09", text: "September" },
      { value: "10", text: "October" },
      { value: "11", text: "November" },
      { value: "12", text: "December" },
    ] as MonthOption[],
    days: Array.from({ length: 31 }, (_, i) => {
      const d = (i + 1).toString().padStart(2, "0");
      return { value: d, text: d };
    }) as DayOption[],
    years: (() => {
      const currentYear = new Date().getFullYear();
      return Array.from({ length: 100 }, (_, i) => {
        const year = (currentYear - i).toString();
        return { value: year, text: year };
      });
    })() as YearOption[],
    captchaImg: "",
    captchaAnswer: "",
  });

  // Handle file uploads
  function handleFileUpload(
    event: Event,
    type: "profilePicture" | "banner"
  ): void {
    const target = event.target as HTMLInputElement;
    const file = target.files?.[0];
    if (!file) return;

    // Validate file type
    const validTypes: string[] = ["image/jpeg", "image/png", "image/gif"];
    if (!validTypes.includes(file.type)) {
      state.errors[type] =
        "Please upload a valid image file (JPEG, PNG, or GIF)";
      return;
    }

    // Validate file size (max 5MB)
    if (file.size > 5 * 1024 * 1024) {
      state.errors[type] = "File size should not exceed 5MB";
      return;
    }

    // Clear previous error
    state.errors[type] = "";

    // Set file and create preview
    state.formData[type] = file;
    const reader = new FileReader();
    reader.onload = (e: ProgressEvent<FileReader>) => {
      if (e.target?.result) {
        state.formData[`${type}Preview`] = e.target.result as string;
      }
    };
    reader.readAsDataURL(file);
  }

  // Upload images to Supabase and get URLs
  async function uploadImagesToSupabase(): Promise<boolean> {
    state.uploadingImages = true;

    try {
      // Upload profile picture if exists
      if (state.formData.profilePicture) {
        const profileUrl = await uploadProfilePicture(
          state.formData.profilePicture
        );

        if (!profileUrl) {
          state.errors.profilePicture = "Failed to upload profile picture";
          state.uploadingImages = false;
          return false;
        }

        state.profilePictureUrl = profileUrl;
      }

      // Upload banner if exists
      if (state.formData.banner) {
        const bannerUrl = await uploadBannerImage(state.formData.banner);

        if (!bannerUrl) {
          state.errors.banner = "Failed to upload banner image";
          state.uploadingImages = false;
          return false;
        }

        state.bannerUrl = bannerUrl;
      }

      state.uploadingImages = false;
      return true;
    } catch (error) {
      // console.error("Error uploading images:", error);
      state.errors.general = "Error uploading images. Please try again.";
      state.uploadingImages = false;
      return false;
    }
  }

  // Prepare data for API submission
  function prepareApiData(): ApiFormData {
    // Format date of birth
    const dob = `${state.formData.birthYear}-${state.formData.birthMonth}-${state.formData.birthDay}`;

    return {
      name: state.formData.name,
      username: state.formData.username,
      email: state.formData.email,
      password: state.formData.password,
      gender: state.formData.gender,
      profilePicture: state.profilePictureUrl || "",
      banner: state.bannerUrl || "",
      dateOfBirth: dob,
      securityQuestion:
        securityQuestions.find(
          (q) => q.id.toString() === state.formData.securityQuestion
        )?.text || "",
      securityAnswer: state.formData.securityAnswer,
      isSubscribed: state.formData.subscribeToNewsletter,
    };
  }

  // Handle step 1 submission
  async function handleStep1Submit(): Promise<void> {
    // console.log("heading to next step")
    state.isLoading = true;
    state.errors.general = "";

    // Validate using external functions
    state.errors.name = validateName(state.formData);
    state.errors.username = await validateUsername(state.formData);
    state.errors.email = await validateEmail(state.formData);
    const pwdRes = validatePassword(state.formData);
    state.errors.password = pwdRes.error;
    state.passwordValidations = pwdRes.validations;
    state.errors.confirmPassword = validateConfirmPassword(state.formData);
    state.errors.gender = validateGender(state.formData);
    state.errors.profilePicture = validateProfilePicture(state.formData);
    state.errors.banner = validateBanner(state.formData);
    state.errors.dateOfBirth = validateDateOfBirth(state.formData);
    const sec = validateSecurityQuestion(state.formData);
    state.errors.securityQuestion = sec.question;
    state.errors.securityAnswer = sec.answer;

    const stepValid =
      !state.errors.name &&
      !state.errors.username &&
      !state.errors.email &&
      !state.errors.password &&
      !state.errors.confirmPassword &&
      !state.errors.gender &&
      !state.errors.profilePicture &&
      !state.errors.banner &&
      !state.errors.dateOfBirth &&
      !state.errors.securityQuestion &&
      !state.errors.securityAnswer &&
      state.recaptchaVerified;

    console.log(stepValid);
    if (stepValid) {
      console.log("heading to next step");

      await sendVerificationCode();

      state.currentStep = 2;
    }
    state.isLoading = false;
  }

  async function sendVerificationCode(): Promise<void> {
    const apiData = {
      email: state.formData.email
    };

    const res = await fetch("http://localhost:8080/api/auth/send-code", {
      method: "POST",
      body: JSON.stringify(apiData)
    });
    if(res.ok) startResendCountdown();
  }

  function startResendCountdown(): void {
    state.resendDisabled = true;
    state.resendCountdown = 60;
    const interval = setInterval(() => {
      state.resendCountdown--;
      if (state.resendCountdown <= 0) {
        clearInterval(interval);
        state.resendDisabled = false;
      }
    }, 1000);
  }

  async function handleResendCode(): Promise<void> {
    if (state.resendDisabled) return;
    state.isLoading = true;

    try {
      await sendVerificationCode();
    } catch (error) {
      state.errors.general =
        "Failed to resend verification code. Please try again.";
    } finally {
      state.isLoading = false;
    }
  }

  async function handleStep2Submit(): Promise<void> {
    state.errors.general = "";
    state.errors.verificationCode = validateVerificationCode(state.formData);
    if (!state.errors.verificationCode) {
      const isVerified = await verifyCode();
      console.log("code verify: " + isVerified);
      if (isVerified) {
        await register();
        navigate("/login");
      } else {
        state.errors.verificationCode = "Invalid verification code";
      }
    }
  }

  async function verifyCode(): Promise<boolean> {
    const apiData = {
      email: state.formData.email,
      code: state.formData.verificationCode
    };

    try {
      const res = await fetch("http://localhost:8080/api/auth/verify-code", {
        method: "POST",
        body: JSON.stringify(apiData)
      });

      const data = await res.json();
      return data.verified;
    } catch (error) {
      console.error("Error verifying code:", error);
      return false;
    };
  }

  // async function completeRegistration(): Promise<void> {
  // }

  function handleGoogleAuth(): void {
    console.log("Google authentication initiated");
  }

  function onRecaptchaCallback(token: string): void {
    fetch("http://localhost:8080/verify-recaptcha", {
      method: "POST",
      // headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ token })
    })
      .then(res => res.json())
      .then(data => {
        if (data.success) {
          console.log("success recaptcha");
          state.recaptchaVerified = true;
        } else {
          state.errors.general = "reCAPTCHA verification failed";
        }
      })
      .catch(() => {
        state.errors.general = "Error verifying reCAPTCHA";
      });
  }

  onMount(async() => {
    console.log("Component mounted, reCAPTCHA would be initialized here");
    // Assign the global callback used by reCAPTCHA v2
    window.recaptchaCallback = onRecaptchaCallback;
    // Dynamically load the reCAPTCHA script if not already on the page
    if (!document.querySelector("script[src=\"https://www.google.com/recaptcha/api.js\"]")) {
      const script = document.createElement("script");
      script.src = "https://www.google.com/recaptcha/api.js";
      script.async = true;
      script.defer = true;
      document.body.appendChild(script);
    }

    await checkGuest(navigate);
  });

  function goBackToStep1(): void {
    state.currentStep = 1;
  }

  function goToLandingPage(): void {
    navigate("/");
  }

  // Testing API endpoint with the new data structure
  async function register() {
    // First upload images to Supabase
    await uploadImagesToSupabase();

    // Then prepare the API data
    const apiData = prepareApiData();
    console.log(JSON.stringify(apiData));

    // Send the request
    const res = await fetch("http://localhost:8080/api/auth/register", {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify(apiData),
    });

    const data = await res.json();
    console.log(data);
  }
</script>

<main class="register-page">
  <div class="register-container">
    <div class="register-header">
      <button class="close-button" type="button" onclick={goToLandingPage}
        >×</button
      >
      <div class="logo">AY</div>
    </div>

    {#if state.currentStep === 1}
      <div class="register-content">
        <h1>Create your account</h1>
        <button
          class="google-auth-button"
          type="button"
          onclick={handleGoogleAuth}
        >
          <span class="google-icon">G</span>
          Sign up with Google
        </button>
        <div class="divider">
          <span>or</span>
        </div>

        <form>
          <FormInput
            id="name"
            type="text"
            label="Name"
            placeholder="Your full name"
            bind:value={state.formData.name}
            error={state.errors.name}
            onBlur={() => (state.errors.name = validateName(state.formData))}
            maxlength={50}
            showCharCount={true}
          />
          <FormInput
            id="username"
            type="text"
            label="Username"
            placeholder="@username"
            bind:value={state.formData.username}
            error={state.errors.username}
            onBlur={async () =>
              (state.errors.username = await validateUsername(state.formData))}
          />
          <FormInput
            id="email"
            type="email"
            label="Email"
            placeholder="email@example.com"
            bind:value={state.formData.email}
            error={state.errors.email}
            onBlur={async () =>
              (state.errors.email = await validateEmail(state.formData))}
          />
          <PasswordInput
            id="password"
            label="Password"
            bind:value={state.formData.password}
            placeholder="Create a password"
            error={state.errors.password}
            onInput={() => {
              const pwd = validatePassword(state.formData);
              state.errors.password = pwd.error;
              state.passwordValidations = pwd.validations;
            }}
            showRequirements={true}
            validations={state.passwordValidations}
          />
          <div class="form-group">
            <label for="confirmPassword">Confirm Password</label>
            <input
              type="password"
              id="confirmPassword"
              bind:value={state.formData.confirmPassword}
              onblur={() =>
                (state.errors.confirmPassword = validateConfirmPassword(
                  state.formData
                ))}
              class:error={state.errors.confirmPassword}
              placeholder="Confirm your password"
            />
            {#if state.errors.confirmPassword}
              <p class="error-message">{state.errors.confirmPassword}</p>
            {/if}
          </div>
          <div class="form-group">
            <p>Gender</p>
            <div class="radio-group">
              <label class="radio-label">
                <input
                  type="radio"
                  name="gender"
                  value="male"
                  bind:group={state.formData.gender}
                  onchange={() =>
                    (state.errors.gender = validateGender(state.formData))}
                />
                Male
              </label>
              <label class="radio-label">
                <input
                  type="radio"
                  name="gender"
                  value="female"
                  bind:group={state.formData.gender}
                  onchange={() =>
                    (state.errors.gender = validateGender(state.formData))}
                />
                Female
              </label>
            </div>
            {#if state.errors.gender}
              <p class="error-message">{state.errors.gender}</p>
            {/if}
          </div>
          <FileUpload
            id="profilePicture"
            label="Profile Picture"
            fileName={state.formData.profilePicture
              ? (state.formData.profilePicture as File).name
              : "Choose file"}
            previewUrl={state.formData.profilePicturePreview}
            error={state.errors.profilePicture}
            onFileChange={(e) => handleFileUpload(e, "profilePicture")}
          />
          <FileUpload
            id="banner"
            label="Banner Image"
            fileName={state.formData.banner
              ? (state.formData.banner as File).name
              : "Choose file"}
            previewUrl={state.formData.bannerPreview}
            error={state.errors.banner}
            isBanner={true}
            onFileChange={(e) => handleFileUpload(e, "banner")}
          />
          <div class="form-group">
            <p>Date of birth</p>
            <p class="dob-info">
              This will not be shown publicly. Confirm your age.
            </p>
            <div class="dob-selects">
              <div class="select-wrapper">
                <select
                  bind:value={state.formData.birthMonth}
                  onchange={() =>
                    (state.errors.dateOfBirth = validateDateOfBirth(
                      state.formData
                    ))}
                >
                  <option value="" disabled selected>Month</option>
                  {#each state.months as month (month.value)}
                    <option value={month.value}>{month.text}</option>
                  {/each}
                </select>
              </div>
              <div class="select-wrapper">
                <select
                  bind:value={state.formData.birthDay}
                  onchange={() =>
                    (state.errors.dateOfBirth = validateDateOfBirth(
                      state.formData
                    ))}
                >
                  <option value="" disabled selected>Day</option>
                  {#each state.days as day (day.value)}
                    <option value={day.value}>{day.text}</option>
                  {/each}
                </select>
              </div>
              <div class="select-wrapper">
                <select
                  bind:value={state.formData.birthYear}
                  onchange={() =>
                    (state.errors.dateOfBirth = validateDateOfBirth(
                      state.formData
                    ))}
                >
                  <option value="" disabled selected>Year</option>
                  {#each state.years as year (year.value)}
                    <option value={year.value}>{year.text}</option>
                  {/each}
                </select>
              </div>
            </div>
            {#if state.errors.dateOfBirth}
              <p class="error-message">{state.errors.dateOfBirth}</p>
            {/if}
          </div>
          <div class="form-group">
            <label for="securityQuestion">Security Question</label>
            <div class="select-wrapper full-width">
              <select
                id="securityQuestion"
                bind:value={state.formData.securityQuestion}
                onchange={() => {
                  const sec = validateSecurityQuestion(state.formData);
                  state.errors.securityQuestion = sec.question;
                  state.errors.securityAnswer = sec.answer;
                }}
              >
                <option value="" disabled selected
                  >Select a security question</option
                >
                {#each securityQuestions as question (question.id)}
                  <option value={question.id.toString()}>{question.text}</option>
                {/each}
              </select>
            </div>
            {#if state.errors.securityQuestion}
              <p class="error-message">{state.errors.securityQuestion}</p>
            {/if}
          </div>
          {#if state.formData.securityQuestion}
            <FormInput
              id="securityAnswer"
              type="text"
              label="Security Answer"
              placeholder="Your answer"
              bind:value={state.formData.securityAnswer}
              error={state.errors.securityAnswer}
              onInput={() => {
                const sec = validateSecurityQuestion(state.formData);
                state.errors.securityAnswer = sec.answer;
              }}
            />
          {/if}
          <div class="form-group checkbox-group">
            <label class="checkbox-label">
              <input
                type="checkbox"
                bind:checked={state.formData.subscribeToNewsletter}
              />
              Subscribe to newsletter
            </label>
          </div>
          <div class="form-group recaptcha">
            <div
              id="recaptcha-container"
              class="g-recaptcha"
              data-sitekey="6Le5ySYrAAAAADAVGZm_Gug20l8ge7D0IQgwg6f3"
              data-callback="recaptchaCallback">
            </div>
          </div>
          {#if state.errors.general}
            <p class="error-message general-error">{state.errors.general}</p>
          {/if}
          <div class="form-actions">
            <Button
              text="Next"
              type="button"
              disabled={state.isLoading}
              onClick={handleStep1Submit}
            />
          </div>
        </form>
      </div>
    {:else if state.currentStep === 2}
      <div class="register-content">
        <h1>We sent you a code</h1>
        <p class="verification-info">
          Enter 5 digit code to verify {state.formData.email}
        </p>
        <form>
          <div class="form-group">
            <label for="verificationCode">Verification code</label>
            <input
              type="text"
              id="verificationCode"
              bind:value={state.formData.verificationCode}
              oninput={() =>
                (state.errors.verificationCode = validateVerificationCode(
                  state.formData
                ))}
              class:error={state.errors.verificationCode}
              placeholder="Enter 5 digit code"
              maxlength="5"
              pattern="[0-9]*"
            />
            {#if state.errors.verificationCode}
              <p class="error-message">{state.errors.verificationCode}</p>
            {/if}
          </div>
          <div class="resend-code">
            <button
              type="button"
              class="text-button"
              onclick={handleResendCode}
              disabled={state.resendDisabled}
            >
              {state.resendDisabled
                ? `Resend code in ${state.resendCountdown}s`
                : "Didn't receive a code?"}
            </button>
          </div>
          {#if state.errors.general}
            <p class="error-message general-error">{state.errors.general}</p>
          {/if}
          <div class="form-actions">
            <button
              type="button"
              class="secondary-button"
              onclick={goBackToStep1}>Back</button
            >
            <button
              type="button"
              class="primary-button"
              disabled={state.isLoading}
              onclick={handleStep2Submit}
            >
              {state.isLoading ? "Verifying..." : "Next"}
            </button>
          </div>
        </form>
      </div>
    {/if}
  </div>
</main>

<style>
  :global(html),
  :global(body) {
    margin: 0;
    padding: 0;
    overflow-x: hidden;
    width: 100%;
    height: 100%;
    background-color: #15202b;
  }

  .register-page {
    background-color: #15202b;
    color: #fff;
    min-height: 100vh;
    width: 100%;
    display: flex;
    align-items: center;
    justify-content: center;
    font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, Oxygen,
      Ubuntu, Cantarell, "Open Sans", "Helvetica Neue", sans-serif;
    padding: 20px;
  }

  .register-container {
    background-color: #000;
    border-radius: 16px;
    width: 100%;
    max-width: 600px;
    overflow: hidden;
    box-shadow: 0 8px 16px rgba(0, 0, 0, 0.2);
  }

  .register-header {
    padding: 16px;
    display: flex;
    align-items: center;
    justify-content: space-between;
    border-bottom: 1px solid #2f3336;
  }

  .close-button {
    background: none;
    border: none;
    color: #fff;
    font-size: 24px;
    cursor: pointer;
    padding: 0;
    margin: 0;
    width: 32px;
    height: 32px;
    display: flex;
    align-items: center;
    justify-content: center;
    border-radius: 50%;
  }

  .close-button:hover {
    background-color: rgba(255, 255, 255, 0.1);
  }

  .logo {
    font-size: 28px;
    font-weight: 900;
    color: #fff;
  }

  .register-content {
    padding: 24px;
  }

  h1 {
    font-size: 24px;
    font-weight: 700;
    margin: 0 0 24px 0;
  }

  .google-auth-button {
    width: 100%;
    padding: 12px;
    border-radius: 9999px;
    border: 1px solid #2f3336;
    background-color: transparent;
    color: #fff;
    font-size: 15px;
    font-weight: 600;
    cursor: pointer;
    display: flex;
    align-items: center;
    justify-content: center;
    margin-bottom: 16px;
    transition: background-color 0.2s;
  }

  .google-auth-button:hover {
    background-color: rgba(255, 255, 255, 0.1);
  }

  .google-icon {
    background-color: #fff;
    color: #000;
    width: 20px;
    height: 20px;
    border-radius: 50%;
    display: flex;
    align-items: center;
    justify-content: center;
    margin-right: 8px;
    font-weight: bold;
  }

  .divider {
    display: flex;
    align-items: center;
    margin: 16px 0;
    color: #6e767d;
  }

  .divider::before,
  .divider::after {
    content: "";
    flex: 1;
    border-bottom: 1px solid #2f3336;
  }

  .divider span {
    padding: 0 10px;
    font-size: 14px;
  }

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

  p {
    display: block;
    font-size: 14px;
    font-weight: 500;
    margin-bottom: 8px;
    color: #e7e9ea;
  }

  input[type="password"] {
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

  .general-error {
    margin-bottom: 16px;
  }

  .radio-group {
    display: flex;
    gap: 16px;
  }

  .radio-label {
    display: flex;
    align-items: center;
    cursor: pointer;
  }

  .radio-label input {
    margin-right: 8px;
  }

  .dob-info {
    font-size: 13px;
    color: #6e767d;
    margin-top: 0;
    margin-bottom: 12px;
  }

  .dob-selects {
    display: flex;
    gap: 12px;
  }

  .select-wrapper {
    position: relative;
    flex: 1;
  }

  .full-width {
    width: 100%;
  }

  select {
    width: 100%;
    padding: 12px;
    border-radius: 4px;
    border: 1px solid #2f3336;
    background-color: #000;
    color: #fff;
    font-size: 16px;
    appearance: none;
    cursor: pointer;
  }

  .select-wrapper::after {
    content: "▼";
    font-size: 12px;
    color: #6e767d;
    position: absolute;
    right: 12px;
    top: 50%;
    transform: translateY(-50%);
    pointer-events: none;
  }

  .checkbox-group {
    display: flex;
    align-items: center;
  }

  .checkbox-label {
    display: flex;
    align-items: center;
    cursor: pointer;
  }

  .checkbox-label input {
    margin-right: 8px;
  }

  .recaptcha {
    margin-bottom: 24px;
  }

  .form-actions {
    display: flex;
    justify-content: space-between;
    gap: 12px;
  }

  .primary-button,
  .secondary-button {
    flex: 1;
    padding: 12px;
    border-radius: 9999px;
    font-size: 15px;
    font-weight: 600;
    cursor: pointer;
    transition: background-color 0.2s;
  }

  .primary-button {
    background-color: #1d9bf0;
    color: white;
    border: none;
  }

  .primary-button:hover:not(:disabled) {
    background-color: #1a8cd8;
  }

  .primary-button:disabled {
    opacity: 0.5;
    cursor: not-allowed;
  }

  .secondary-button {
    background-color: transparent;
    color: #fff;
    border: 1px solid #2f3336;
  }

  .secondary-button:hover {
    background-color: rgba(255, 255, 255, 0.1);
  }

  .verification-info {
    font-size: 15px;
    color: #6e767d;
    margin-bottom: 24px;
  }

  .text-button {
    background: none;
    border: none;
    color: #1d9bf0;
    font-size: 14px;
    cursor: pointer;
    padding: 0;
  }

  .text-button:hover:not(:disabled) {
    text-decoration: underline;
  }

  .text-button:disabled {
    opacity: 0.5;
    cursor: not-allowed;
  }

  .resend-code {
    margin-bottom: 24px;
  }

  /* Responsive styles */
  @media (max-width: 480px) {
    .register-container {
      max-width: 100%;
      height: 100vh;
      border-radius: 0;
    }

    .dob-selects {
      flex-direction: column;
      gap: 8px;
    }

    .form-actions {
      flex-direction: column;
    }
  }
</style>
