import { create } from 'zustand'

type State = {
  themeMode: boolean
}

type Action = {
  updateThemeMode: (themeMode: State['themeMode']) => void
}


const useThemeMode = create<State & Action>((set) => ({
  themeMode: false,
  updateThemeMode: (themeMode:boolean) => set(() => ({ themeMode: themeMode })),
}))

export default useThemeMode


//use it like this in component:
// const [themeMode, updateThemeMode] = useThemeMode(
//   (state:any) => [state.themeMode, state.updateThemeMode],
// )


// useEffect(() => {
//   if (window) {
//     let html = window.document.documentElement.classList;
//     if (html[0]==="dark") {
//       updateThemeMode(true)
//     }else{
//       updateThemeMode(false)
//     }
//   }  
  

//   console.log(themeMode,"themeMode")
// }, [themeMode, updateThemeMode]);