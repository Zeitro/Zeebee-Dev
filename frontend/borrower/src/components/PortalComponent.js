'use client'
import React, {
  useState,
  useRef,
  useEffect,
} from "react";
import ReactDOM from "react-dom";
import eventBus from "@/lib/bus";

/**
 * PortalComponent是将子组件渲染到body下id为portal-root的div中
 */
const PortalComponent = () => {
  const [content, setContent] = useState("");
  const portalRootRef = useRef(null);

  const onClose = () => {
    if (portalRootRef.current) {
      ReactDOM.unmountComponentAtNode(portalRootRef.current);
      document.getElementById("portal-root") &&document.body.removeChild(portalRootRef.current);
      portalRootRef.current = null;
      setContent('');
    }
  };

  const appendPortalComponent = (element) => {
    let portalEle = document.getElementById("portal-root");
    if (!portalEle) {
      portalEle = document.createElement("div");
      portalEle.id = "portal-root";
      document.body.appendChild(portalEle);
    }
    portalRootRef.current = portalEle;
    setContent(ReactDOM.createPortal(element, portalRootRef.current));
  };


  useEffect(() => {
    eventBus?.on("appendPortalComponent", appendPortalComponent);
    eventBus?.on("removePortalComponent", onClose);
    return () => {
        eventBus?.removeListener("appendPortalComponent", appendPortalComponent);
        eventBus?.removeListener("removePortalComponent", onClose);
        onClose();
    };
  }, []);

  // 只有 portalRootRef.current 存在时才渲染 portal
  if (!portalRootRef.current) return null;

  return content;
};

export default PortalComponent;
