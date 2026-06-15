// 导入 react-pdf 所需的 CSS
import '@/app/loans/shared/importReactPdfStyles';

// 这个布局组件接收并渲染它的子组件（也就是你的 page.tsx）
export default function ContractLayout({ children }) {
	return <>{children}</>;
}
