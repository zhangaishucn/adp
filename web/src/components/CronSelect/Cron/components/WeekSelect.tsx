// @ts-nocheck
import { useEffect } from 'react';
import intl from 'react-intl-universal';
import { Select } from 'antd';
import locales from '../../locales';

const weekOptions = {
  1: intl.get('CronSelect.SUN'),
  2: intl.get('CronSelect.MON'),
  3: intl.get('CronSelect.TUE'),
  4: intl.get('CronSelect.WED'),
  5: intl.get('CronSelect.THU'),
  6: intl.get('CronSelect.FRI'),
  7: intl.get('CronSelect.SAT'),
};

export { weekOptions as weekOptionsObj };

function WeekSelect(props): JSX.Element {
  useEffect(() => {
    intl.load(locales);
  }, []);

  return (
    <Select size="small" {...props}>
      {Object.entries(weekOptions).map(([weekCode, weekName]) => {
        return (
          <Select.Option key={weekCode} value={weekCode}>
            {weekName}
          </Select.Option>
        );
      })}
    </Select>
  );
}

export default WeekSelect;
