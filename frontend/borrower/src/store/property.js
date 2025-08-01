import { create } from "zustand"
export const usePropertyStore = create((set) => ({
  propertyTypeOptions: [],
  login: (userData) => set({ userInfo: userData })
}))