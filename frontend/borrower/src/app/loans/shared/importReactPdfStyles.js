/**
 * react-pdf 所需全局样式。
 * 正常走路由时由 createZeeBeeLoan / zeeBeeDetails 的 layout.js 引入；
 * Tab 保活经 tabRegistry 懒加载 page 时需经 tabEntry 或 layout 一并引入。
 */
import 'react-pdf/dist/Page/AnnotationLayer.css';
import 'react-pdf/dist/Page/TextLayer.css';
