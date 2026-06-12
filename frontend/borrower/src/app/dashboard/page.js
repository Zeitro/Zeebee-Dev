'use client';
import {message, Radio, Select, Spin} from 'antd';
import React, {useCallback, useEffect, useMemo, useRef, useState} from 'react';
import MaturityAlerts from './components/MaturityAlerts';
import LoanStatusCard from './components/LoanStatusCard';
import DifferentStatusLoans from './components/loans';
import {FieldTips} from '@repo/shared-methods';
import {filterTimeCount, filterTimeCountTips, loansSourceTypeOptions,} from '@repo/shared-methods/enum';
import {useUserStore} from '@repo/shared-methods/store';
import {createRequest, getDateRange, sessionStorageUtil} from '@repo/shared-methods/utils';
import {useAuth} from '@repo/shared-methods/hooks';
import {PERMISSIONS} from '@config/permissions';
import {useRouter} from 'next/navigation';
// 授权gmail，以用户邮箱发送邮件<中转页,接收参数>
const Dashboard = () => {
    const router = useRouter();
    const {hasPermission} = useAuth();
    const {request} = createRequest();
    const {userInfo} = useUserStore();
    const [type, setType] = useState('All');
    const [dayCount, setDayCount] = useState('All');
    const [data, setData] = useState({});
    const [loading, setLoading] = useState(true);
    const headerRef = useRef();
    // 区分哪类loans
    const source = useMemo(() => {
        let v = [];
        if (type === 'ZeeBee') {
            v = ['System', 'Preapproval'];
        } else if (type === '3rdParty') {
            v = ['Borrower'];
        }
        return v;
    }, [type]);

    // 获取时间段startDate与endDate
    const timePeriod = useMemo(() => {
        return getDateRange(dayCount);
    }, [dayCount]);

    const getData = useCallback(async () => {
        try {
            setLoading(true);

            const res = await request('/loan/dashboard', {
                method: 'POST',
                body: {
                    ...timePeriod,
                    source,
                },
            });
            setData(res);
            setLoading(false);
        } catch (error) {
            setLoading(false);
            throw error;
        }
    }, [source, timePeriod]);
    const handleChange = (e) => {
        setType(e.target.value);
    };

    const handleGotoLoans = (stage, date_type) => {
        if (!hasPermission(PERMISSIONS?.LOAN?.LOAN_PAGE)) {
            message.open({
                content:
                    'You do not have permission to access this page. Please contact the administrator to grant you access.',
                type: 'warning',
            });
            return;
        }
        router.push('/loans');
        sessionStorageUtil.set('LOANS-FILTER', {
            stage,
            sourceType: type,
            dayCount,
            sorter: {
                column: date_type,
                sort: 'desc',
            },
            date_type,
        });
    };

    const handleGotoLoanDetails = (item) => {
        if (!hasPermission(PERMISSIONS?.LOAN?.DETAILS)) {
            message.open({
                content:
                    'You do not have permission to access this page. Please contact the administrator to grant you access.',
                type: 'warning',
            });
            return;
        }
        let url = `/loans/zeeBeeDetails?id=${item.id}`;
        if (item.source === 'Borrower') {
            url = `/loans/${item.id}`;
        }
        router.push(url);
    };

    useEffect(() => {
        const initData = async () => {
            await getData();
        };
        initData();
    }, [getData, type, dayCount]);
    return (
        <div className="h-full w-full overflow-hidden overflow-y-auto px-6 pb-6 bg-[#f5f5f5]">
            <Spin spinning={loading}>
                <div
                    ref={headerRef}
                    className={`flex items-center justify-between flex-col lg:flex-row pb-6 pt-7 sticky top-0 bg-[#f5f5f5] z-20
								`}
                >
					<span className="text-3xl/10 font-medium text-base-text pb-2 lg:pb-0">
						Hi {userInfo?.first_name},welcome back!
					</span>
                    <Radio.Group
                        onChange={handleChange}
                        value={type}
                        options={loansSourceTypeOptions}
                        optionType="button"
                        block
                        className="w-[380px]"
                    />
                </div>
                <Select
                    className="min-w-[150px] !mb-4"
                    onChange={(v) => setDayCount(v)}
                    value={dayCount}
                    options={filterTimeCount}
                    suffixIcon={<FieldTips title={filterTimeCountTips} isShowDescription={true}/>}
                />
                <div className="grid grid-cols-1 lg:grid-cols-3 gap-2">
                    <LoanStatusCard
                        dayCount={dayCount}
                        data={data}
                        timePeriod={timePeriod}
                        handleGotoLoans={handleGotoLoans}
                    />
                    <MaturityAlerts
                        data={data}
                        source={source}
                        handleGotoLoanDetails={handleGotoLoanDetails}
                    />
                </div>
                <DifferentStatusLoans
                    stage={'Pending'}
                    handleGotoLoans={handleGotoLoans}
                    handleGotoLoanDetails={handleGotoLoanDetails}
                    data={data}
                    dataType={'closing_date'}
                />
                <DifferentStatusLoans
                    stage={'Active'}
                    handleGotoLoans={handleGotoLoans}
                    handleGotoLoanDetails={handleGotoLoanDetails}
                    data={data}
                    dataType={'updated_at'}
                />
            </Spin>
        </div>
    );
};
export default Dashboard;
