// 导入 react-pdf 所需的 CSS
import 'react-pdf/dist/Page/AnnotationLayer.css';
import 'react-pdf/dist/Page/TextLayer.css';

// 这个布局组件接收并渲染它的子组件（也就是你的 page.tsx）
export default function ContractLayout({ children }) {
	return <>{children}</>;
}
