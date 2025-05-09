export type SecurityQuestion = {
    id: number;
    text: string;
};

export type MonthOption = {
    value: string;
    text: string;
};

export type DayOption = {
    value: string;
    text: string;
};

export type YearOption = {
    value: string;
    text: string;
};

export type FormData = {
    name: string;
    username: string;
    email: string;
    password: string;
    confirmPassword: string;
    gender: string;
    profilePicture: File | null;
    profilePicturePreview: string;
    banner: File | null;
    bannerPreview: string;
    birthDay: string;
    birthMonth: string;
    birthYear: string;
    securityQuestion: string;
    securityAnswer: string;
    captchaAnswer: string;
    subscribeToNewsletter: boolean;
    verificationCode: string;
};

export type ErrorMessages = {
    name: string;
    username: string;
    email: string;
    password: string;
    confirmPassword: string;
    gender: string;
    profilePicture: string;
    banner: string;
    dateOfBirth: string;
    securityQuestion: string;
    securityAnswer: string;
    captchaAnswer: string;
    verificationCode: string;
    general: string;
};

export type PasswordValidations = {
    length: boolean;
    uppercase: boolean;
    lowercase: boolean;
    number: boolean;
};

export function validateName(formData: FormData): string {
    const nameRegex = /^[a-zA-Z\s]{4,}$/;
    if (!formData.name) return "Name is required";
    if (!nameRegex.test(formData.name))
        return "Name must be at least 4 characters and contain only letters";
    return "";
}

export async function validateUsername(formData: FormData): Promise<string> {
    if (!formData.username) return "Username is required";
    if (formData.username.length < 4) return "Username must be at least 4 characters";
    return "";
}

export async function validateEmail(formData: FormData): Promise<string> {
    const emailRegex = /^[a-zA-Z0-9._-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,6}$/;
    if (!formData.email) return "Email is required";
    if (!emailRegex.test(formData.email))
        return "Please enter a valid email address";
    return "";
}

export function validatePassword(formData: FormData): { error: string; validations: PasswordValidations } {
    const validations: PasswordValidations = {
        length: formData.password.length >= 8,
        uppercase: /[A-Z]/.test(formData.password),
        lowercase: /[a-z]/.test(formData.password),
        number: /[0-9]/.test(formData.password),
    };

    if (!formData.password) return { error: "Password is required", validations };
    if (!(validations.length && validations.uppercase && validations.lowercase && validations.number))
        return { error: "Password does not meet requirements", validations };
    return { error: "", validations };
}

export function validateConfirmPassword(formData: FormData): string {
    if (!formData.confirmPassword) return "Please confirm your password";
    if (formData.password !== formData.confirmPassword) return "Passwords do not match";
    return "";
}

export function validateGender(formData: FormData): string {
    return formData.gender ? "" : "Please select your gender";
}

export function validateDateOfBirth(formData: FormData): string {
    if (!formData.birthDay || !formData.birthMonth || !formData.birthYear)
        return "Please select your date of birth";
    const birthDate = new Date(
        parseInt(formData.birthYear),
        parseInt(formData.birthMonth) - 1,
        parseInt(formData.birthDay)
    );
    const today = new Date();
    let age = today.getFullYear() - birthDate.getFullYear();
    const monthDiff = today.getMonth() - birthDate.getMonth();
    if (monthDiff < 0 || (monthDiff === 0 && today.getDate() < birthDate.getDate())) {
        age--;
    }
    return age >= 13 ? "" : "You must be at least 13 years old to register";
}

export function validateSecurityQuestion(formData: FormData): { question: string; answer: string } {
    const questionError = formData.securityQuestion ? "" : "Please select a security question";
    const answerError = formData.securityAnswer ? "" : "Please provide an answer";
    return { question: questionError, answer: answerError };
}

export function validateProfilePicture(formData: FormData): string {
    return formData.profilePicture ? "" : "Please upload a profile picture";
}

export function validateBanner(formData: FormData): string {
    return formData.banner ? "" : "Please upload a banner image";
}

export function validateVerificationCode(formData: FormData): string {
    if (!formData.verificationCode) return "Please enter the verification code";
    if (formData.verificationCode.length !== 5) return "Verification code must be 5 digits";
    return "";
}