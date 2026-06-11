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

export const useAuthStore = defineStore('auth', () => {
  const user = ref<User | null>(null)
  const isAdmin = ref(false)
  const loading = ref(true)

  // This promise will resolve once the first auth state check is complete
  const isReady = new Promise((resolve) => {
    onAuthStateChanged(auth, async (firebaseUser) => {
      user.value = firebaseUser
      if (firebaseUser) {
        // Check for admin custom claim
        const token = await firebaseUser.getIdTokenResult()
        isAdmin.value = !!token.claims.admin
      } else {
        isAdmin.value = false
      }
      loading.value = false
      resolve(true)
    })
  })

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

  return { user, isAdmin, loading, isReady, loginWithGoogle, signUp, loginWithEmail, logout };
});
