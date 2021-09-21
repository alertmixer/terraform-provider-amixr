To use provider locally:

1. Clone provider repository:
   
   ```bash
   git clone git@github.com:alertmixer/terraform-provider-amixr
   ```
   
2. Checkout branch with actual provider version:
   
   ```bash
   git checkout schedules_v2
   ```
   
3. If it is needed to change amixr backend url - go to the amixr/config.go and change line 12:

   From:
   ```
   client, err := amixr.NewClient(c.Token)
   ```
   
   To:
   ```
   client, err := amixr.NewClientWithCustomUrl(c.Token, "url_you_want_to_use_locally")
   ```

4. Run:
   
   ```bash
   make install
   ```

5. Import:

   ```hcl
   terraform {
      required_providers {
         amixr = {
            source = "amixr.io/alertmixer/amixr"
            version = "0.2.3"
         }
      }
   }
   ```
