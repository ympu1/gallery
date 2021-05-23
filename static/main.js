$(document).ready(function()
{
	var files_array = [];
	var i = 0;
	$('.image, .video').each(function()
	{
		var file  = {}
		file.position = i;
		i++;
		file.url  = $(this).attr('href');
		if ($(this).hasClass('video'))
		{
			file.type = "video";
		}
		else
		{
			file.type = "image"
		}

		files_array.push(file);

		$(this).click(function()
		{
			change_slide(file);
			$('body').addClass('popup_is_showed')
			$('#popup').fadeIn();
			return false;
		});
	});

	$('#popup').click(function(e)
	{
		var current_slide_position = $(this).data('position');
		var cursor_point_percent = e.pageX / window.innerWidth * 100;
		if (cursor_point_percent >= 50)
		{
			show_next_slide(current_slide_position);
		}
		else
		{
			show_prev_slide(current_slide_position);
		}
		return false;
	});

	function change_slide(file)
	{
		var popup = $('#popup');
		if (file.type == 'image')
		{
			popup.html('');
			popup.css('background-image', 'url("' + file.url + '")');
		}
		else if (file.type == 'video')
		{
			popup.css('background-image', 'none');
			popup.html('<video src="' + file.url + '" loop muted preload="auto"></video>').trigger('play');
			popup.find('video').trigger('play');
		}

		popup.data('position', file.position);
	}

	function show_prev_slide(current_position)
	{
		var prev_position = current_position - 1;
		if (!(prev_position in files_array))
		{
			prev_position = files_array.length - 1;
		}

		var file = files_array[prev_position];
		change_slide(file);
	}

	function show_next_slide(current_position)
	{
		var next_position = current_position + 1;
		if (!(next_position in files_array))
		{
			next_position = 0;
		}

		var file = files_array[next_position];
		change_slide(file);
	}
});