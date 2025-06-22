INSERT INTO users (id, username, password, role)
VALUES (1, 'admin',
        '$argon2id$v=19$m=65536,t=1,p=16$WBBTUigKowG059hx+JpcPg$LRthl+d/F5AWdWIaf2rpEuAVD2woZZ6w99+mBVRdYNY', 'user')
ON CONFLICT DO NOTHING;

INSERT INTO users (id, username, password, role)
VALUES (2, 'user', '$argon2id$v=19$m=65536,t=1,p=16$ch41s7N/Thl8Qo/6fZuTTg$dvoPPS7d5S0WudiPs/pXEX3CYOMjVGkijVEX3V9Ti4w',
        'user')
ON CONFLICT DO NOTHING;
