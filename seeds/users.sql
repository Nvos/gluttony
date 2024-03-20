INSERT INTO public.users (id, name, password)
VALUES (1,
        'admin',
        '$argon2id$v=19$m=19456,t=2,p=1$Yf9sP8D3HeXqyE/tNdTp0Q$VmxkmBZk7WEAw8mHzZrUYoZuK2+mAPsePIKGRktvtsQ')
ON CONFLICT DO NOTHING;
