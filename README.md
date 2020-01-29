# FDFapp

This application is meant to make it easy to report participants for each activity in you local FDF group.

## Database Setup

You will also need to make sure that **you** start/install the database of your choice. Buffalo **won't** install and start it for you.

For example:

	$ docker-compose up

### Create Your Databases

Ok, so you've edited the "database.yml" file and started your database, now Buffalo can create the databases in that file for you:

	$ buffalo pop create -a
	$ buffalo db migrate up

## Starting the Application

Buffalo ships with a command that will watch your application and automatically rebuild the Go binary and any assets for you. To do that run the "buffalo dev" command:

	$ buffalo dev

If you point your browser to [http://127.0.0.1:3000](http://127.0.0.1:3000) you should see a "Welcome to Buffalo!" page.

**Congratulations!** You now have your Buffalo application up and running.

[Powered by Buffalo](http://gobuffalo.io)
