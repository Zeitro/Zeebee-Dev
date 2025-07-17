'use client';
import Image from 'next/image';
import styles from './index.module.scss';
import { useUserStore } from '@/store/user';
import { ConfigProvider, Layout } from 'antd';
import SiderBar from '@/components/Sidebar';
import NewHeader from '@/components/NewHeader';
const { Header, Footer, Sider, Content } = Layout;
export default function Home() {
    const { userInfo } = useUserStore();
    console.log(userInfo);

    return (
       <div>
			11666
		</div>
    );
}
