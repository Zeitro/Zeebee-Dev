import React, { useState, forwardRef, useImperativeHandle, useRef, useLayoutEffect } from "react";
import ReactDOM from "react-dom";
/**
 * PortalComponent是将子组件渲染到body下id为portal-root的div中
 */
const PortalComponent = forwardRef(({ children }, ref) => {
    const [mounted, setMounted] = useState(false);
    const portalRootRef = useRef(null);

    useLayoutEffect(() => {
        const div = document.createElement("div");
        div.id = "portal-root";
        document.body.appendChild(div);
        portalRootRef.current = div;
        setMounted(true); // 触发重新渲染
        return () => {
            ReactDOM.unmountComponentAtNode(div);
            document.getElementById("portal-root") && document.body.removeChild(div);
            portalRootRef.current = null;
            setMounted(false);
        };
    }, []);

    useImperativeHandle(ref, () => ({
        onClose: () => {
            if (portalRootRef.current) {
                ReactDOM.unmountComponentAtNode(portalRootRef.current);
                document.body.removeChild(portalRootRef.current);
                portalRootRef.current = null;
                setMounted(false);
            }
        }
    }));

    // 只有 portalRootRef.current 存在时才渲染 portal
    if (!mounted || !portalRootRef.current) return null;

    return ReactDOM.createPortal(
        children,
        portalRootRef.current
    );
});

export default PortalComponent;