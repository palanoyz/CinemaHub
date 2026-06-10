import { defineStore } from 'pinia';
import { ref } from 'vue';
import { auth, googleProvider } from '../firebase';
import {
  signInWithPopup,
  signOut,
  onAuthStateChanged,
  createUserWithEmailAndPassword,
  signInWithEmailAndPassword,
  type User
} from 'firebase/auth';

async function loginWithGoogle() {
  return signInWithPopup(auth, googleProvider);
}

async function signUp(email: string, password: string) {
  return createUserWithEmailAndPassword(auth, email, password);
}

async function loginWithEmail(email: string, password: string) {
  return signInWithEmailAndPassword(auth, email, password);
}

async function logout() {
  return signOut(auth);
}

export const useAuthStore = defineStore('auth', () => {
  const user = ref<User | null>(null);
  const loading = ref(true);

  // This promise will resolve once the first auth state check is complete
  const isReady = new Promise((resolve) => {
    onAuthStateChanged(auth, (firebaseUser) => {
      user.value = firebaseUser;
      loading.value = false;
      resolve(true);
    });
  });

  return { user, loading, isReady, loginWithGoogle, signUp, loginWithEmail, logout };
});
